package main

import "gorm.io/gorm"

type TLDRDBCached struct {
	provider TLDRProvider
	db       *gorm.DB
}

func (t *TLDRDBCached) Retrieve(key string) string {
	var entity TLDREntity
	result := t.db.First(&entity, "key = ?", key)

	if result.Error == nil {
		return entity.Val
	}

	val := t.provider.Retrieve(key)

	newEntity := TLDREntity{
		Key: key,
		Val: val,
	}
	t.db.Create(&newEntity)

	return val
}

func (t *TLDRDBCached) List() []string {
	providerList := t.provider.List()

	var cachedEntities []TLDREntity
	t.db.Find(&cachedEntities)

	cachedKeysMap := make(map[string]bool)
	for _, e := range cachedEntities {
		cachedKeysMap[e.Key] = true
	}

	result := []string{}
	for _, e := range cachedEntities {
		result = append(result, e.Key)
	}

	for _, key := range providerList {
		if !cachedKeysMap[key] {
			result = append(result, key)
		}
	}

	return result
}

func NewTLDRDBCached(nonCachedProvider TLDRProvider) TLDRProvider {
	db := GetConnection()
	db.AutoMigrate(&TLDREntity{})

	return &TLDRDBCached{
		provider: nonCachedProvider,
		db:       db,
	}
}
