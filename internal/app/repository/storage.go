package repository

var mapStorage = map[string]string{}

// Сохранить в мапу оба линка
func SaveLinks(shortLink, link string) {
	mapStorage[shortLink] = link
}

// Вернуть full link
func GetFullLink(shortLink string) string {
	return mapStorage[shortLink]
}

// Проверить, не был ли ранее записан этот шорт линк
func CheckAlreadyHaveShortLink(link string) (bool, string) {
	for k, v := range mapStorage {
		if v == link {
			return true, k
		}
	}		
	return false, ""	
}