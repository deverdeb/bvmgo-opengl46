package assetsmngr

import (
	"fmt"
	"log/slog"
)

// AssetLoadFunction permet de définir une fonction de chargement d'une ressource.
type AssetLoadFunction[T any] func() (*T, error)

// AssetReleaseFunction permet de définir une fonction de libération d'une ressource.
type AssetReleaseFunction[T any] func(*T)

// assetInformation contient les informations permettant de charger / libérer une ressource.
type assetInformation[T any] struct {
	// Name est le nom de la ressource.
	Name string
	// Load est la fonction de chargement de la ressource.
	Load AssetLoadFunction[T]
	// Release est la fonction de libération de la ressource.
	Release AssetReleaseFunction[T]
	// Asset est un pointeur vers la ressource.
	Asset *T
	// IsObsolete est un flag indiquant que la ressource est obsolète
	IsObsolete bool
}

// Manager est un gestionnaire de ressources
type Manager[T any] struct {
	// Ensemble des informations sur les ressources, par nom.
	assetInformations map[string]*assetInformation[T]
}

// CreateManager crée et initialise un gestionnaire de ressources
func CreateManager[T any]() Manager[T] {
	return Manager[T]{
		assetInformations: make(map[string]*assetInformation[T]),
	}
}

// RegisterAsset enregistre une ressource non gérée par le manager.
func (manager *Manager[T]) RegisterAsset(name string, asset *T) {
	oldInformation, found := manager.assetInformations[name]
	if found && oldInformation.Asset == asset {
		oldInformation.IsObsolete = false
	} else {
		manager.registerInformation(&assetInformation[T]{
			Name:       name,
			Load:       nil,
			Release:    nil,
			Asset:      asset,
			IsObsolete: false,
		})
	}
}

// Register enregistre une ressource avec une fonction de chargement et de libération (optionnelle).
func (manager *Manager[T]) Register(name string, loadFunction AssetLoadFunction[T], releaseFunction AssetReleaseFunction[T]) error {
	if loadFunction == nil {
		// Ressource sans méthode de chargement
		return fmt.Errorf("cannot registered '%s' asset, because loading method is missing", name)
	}
	if releaseFunction == nil {
		// Ressource sans méthode de libération
		slog.Warn("cannot registered asset, because release method is missing", "assert name", name)
		return nil
	}
	oldInformation, found := manager.assetInformations[name]
	if found {
		oldInformation.IsObsolete = false
	} else {
		manager.registerInformation(&assetInformation[T]{
			Name:       name,
			Load:       loadFunction,
			Release:    releaseFunction,
			Asset:      nil,
			IsObsolete: false,
		})
	}
	return nil
}

// RegisterInformation enregistre une ressource
func (manager *Manager[T]) registerInformation(assetInformation *assetInformation[T]) {
	// Au cas où la ressource est déjà enregistrée
	_, found := manager.assetInformations[assetInformation.Name]
	if found {
		slog.Warn("asset is already registered", "assert name", assetInformation.Name)
		manager.unregisterInformation(assetInformation)
	}
	// Enregistrer la ressource
	manager.assetInformations[assetInformation.Name] = assetInformation
}

// Unregister désenregistre une ressource du manager (et supprime la ressource associée)
func (manager *Manager[T]) Unregister(name string) {
	assetInformation, found := manager.assetInformations[name]
	if found {
		manager.unregisterInformation(assetInformation)
	}
}

// UnregisterAllObsolete désenregistre toutes les ressources obsolètes (et supprime la ressource associée)
func (manager *Manager[T]) UnregisterAllObsolete() {
	obsoletes := make([]*assetInformation[T], 0, 32)
	for _, assetInformation := range manager.assetInformations {
		if assetInformation.IsObsolete {
			obsoletes = append(obsoletes, assetInformation)
		}
	}
	for _, assetInformation := range obsoletes {
		manager.unregisterInformation(assetInformation)
	}
}

// UnregisterAll désenregistre toutes les ressources (et supprime la ressource associée)
func (manager *Manager[T]) UnregisterAll() {
	for _, assetInformation := range manager.assetInformations {
		manager.unregisterInformation(assetInformation)
	}
}

// unregisterInformation désenregistre une ressource du manager (et supprime la ressource associée)
func (manager *Manager[T]) unregisterInformation(assetInformation *assetInformation[T]) {
	if assetInformation.Asset != nil && assetInformation.Load != nil {
		manager.releaseInformation(assetInformation)
	}
	delete(manager.assetInformations, assetInformation.Name)
}

// UnregisterAsset libère une ressource
func (manager *Manager[T]) UnregisterAsset(asset *T) {
	if asset == nil {
		return
	}
	// Chercher la ressource
	for _, assetInformation := range manager.assetInformations {
		if assetInformation.Asset == asset {
			manager.unregisterInformation(assetInformation)
			return
		}
	}
	slog.Warn("cannot unregistered assert because asset is not found")
}

// LoadAll charge toutes les ressources non chargées
func (manager *Manager[T]) LoadAll() error {
	withError := 0
	var lastError error = nil // FIXME utiliser les groupes d'erreurs de GO 1.21
	// Parcourir toutes les ressources
	for _, assetInformation := range manager.assetInformations {
		if assetInformation.Asset == nil && assetInformation.Load != nil && !assetInformation.IsObsolete {
			err := manager.loadInformation(assetInformation)
			if err != nil {
				withError += 1
				slog.Error("load asset failed", "asset name", assetInformation.Name, "error", err)
				lastError = err
			}
		}
	}
	if withError > 0 {
		return fmt.Errorf("load all resources failed with %d error(s)\n - %w", withError, lastError)
	}
	return nil
}

// loadInformation charge une ressource via ses informations
func (manager *Manager[T]) loadInformation(assetInformation *assetInformation[T]) error {
	// Faire quelques vérifications pré-chargement
	if assetInformation.Asset != nil {
		slog.Info("cannot load asset because it is already loaded", "asset name", assetInformation.Name)
		return nil
	}
	if assetInformation.IsObsolete {
		slog.Warn("cannot load asset because it is obsolete", "asset name", assetInformation.Name)
		return nil
	}
	if assetInformation.Load == nil {
		slog.Warn("cannot load asset because missing load function for it", "asset name", assetInformation.Name)
		return nil
	}
	// Charger la ressource
	asset, err := assetInformation.Load()
	if err != nil {
		return fmt.Errorf("failed load '%s' asset\n - %w", assetInformation.Name, err)
	}
	assetInformation.Asset = asset
	return nil
}

// Release libère une ressource, mais sans la supprimer du manager.
func (manager *Manager[T]) Release(name string) {
	assetInformation, found := manager.assetInformations[name]
	if !found {
		slog.Warn("cannot release asset, because it is not found", "asset name", name)
		return
	}
	manager.unregisterInformation(assetInformation)
}

// ReleaseAll force la libération de toutes les ressources, mais sans les supprimer du manager
func (manager *Manager[T]) ReleaseAll() {
	// Libérer toutes les ressources
	for _, assetInformation := range manager.assetInformations {
		if assetInformation.Load != nil {
			manager.releaseInformation(assetInformation)
		}
	}
}

// releaseInformation libère une ressource (qui reste au niveau du manager et pourra être rechargée)
func (manager *Manager[T]) releaseInformation(assetInformation *assetInformation[T]) {
	if assetInformation.Asset == nil {
		// Ressource non chargée
		return
	}
	if assetInformation.Load == nil {
		// Ressource non managée
		slog.Warn("cannot release asset because it is not managed by manager", "asset name", assetInformation.Name)
		return
	}
	if assetInformation.Release == nil {
		// Ressource sans méthode de libération
		slog.Warn("cannot release asset because release method is not defined", "asset name", assetInformation.Name)
		return
	}
	// Libérer la ressource
	assetInformation.Release(assetInformation.Asset)
	assetInformation.Asset = nil
}

// ReleaseAllObsolete force la libération des ressources obsolètes
func (manager *Manager[T]) ReleaseAllObsolete() {
	// Libérer toutes les ressources
	for _, assetInformation := range manager.assetInformations {
		if assetInformation.IsObsolete && assetInformation.Asset != nil && assetInformation.Load != nil {
			manager.releaseInformation(assetInformation)
		}
	}
}

// ReleaseAsset libère une ressource
func (manager *Manager[T]) ReleaseAsset(asset *T) {
	if asset == nil {
		return
	}
	// Chercher la ressource
	for _, assetInformation := range manager.assetInformations {
		if assetInformation.Asset == asset {
			manager.releaseInformation(assetInformation)
			return
		}
	}
	slog.Warn("asset is not found")
}

// MaskObsolete marque une ressource comme obsolète
func (manager *Manager[T]) MaskObsolete(name string) {
	assetInformation, found := manager.assetInformations[name]
	if found {
		assetInformation.IsObsolete = true
	}
}

// MaskNotObsolete marque une ressource comme non obsolète
func (manager *Manager[T]) MaskNotObsolete(name string) {
	assetInformation, found := manager.assetInformations[name]
	if found {
		assetInformation.IsObsolete = false
	}
}

// Get permet de récupérer une ressource
func (manager *Manager[T]) Get(name string) (*T, error) {
	// Récupérer les informations sur la ressource
	assetInformation, found := manager.assetInformations[name]
	if !found {
		return nil, fmt.Errorf("cannot find '%s' asset because it is not registered", name)
	}
	// Vérifier que la ressource est présente
	if assetInformation.Asset == nil {
		if assetInformation.Load == nil {
			return nil, fmt.Errorf("cannot load '%s' asset because it has not loading method", name)
		}
		// Ressource non présente, la charger
		err := manager.loadInformation(assetInformation)
		if err != nil {
			return nil, err
		}
	}
	// Retourner la ressource liée à l'information
	return assetInformation.Asset, nil
}
