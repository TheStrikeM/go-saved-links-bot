package telegram

const msgHelp = `Я крутой бот, просто используй меня. Без лишних слов.

Слава КПСС. Слава партии. Слава Бауманке.

А вообще могу сохранять ссылки, а потом выдавать рандомную из них /rnd
После команды /rnd удаляю ссылку, аккуратно
`

const msgHello = "Здарова! \n\n" + msgHelp

const prefix = "| StrikeBot | "
const (
	msgUnknownCommand = prefix + "Неизвестная команда"
	msgNoSavedPages   = prefix + "У вас нет сохраненных ссылок"
	msgSaved          = prefix + "Сохранил твою ссылку"
	msgAlreadyExists  = prefix + "Эта ссылка уже существует, друг"
)
