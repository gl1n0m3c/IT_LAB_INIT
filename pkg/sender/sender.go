package sender

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/config"
	"github.com/spf13/viper"
	"net/http"
	"net/smtp"
	"os"
	"strings"
)

var (
	mail     string
	password string
	host     string
	port     string
)

func InitEmailConfig() {
	mail = viper.GetString(config.Mail)
	password = viper.GetString(config.MailPassword)
	host = viper.GetString(config.MailHost)
	port = viper.GetString(config.MailPort)
}

type Sender struct {
	auth smtp.Auth
}

type Message struct {
	To         []string
	Subject    string
	Body       string
	Attachment []byte
}

func New() *Sender {
	auth := smtp.PlainAuth("", mail, password, host)
	return &Sender{auth}
}

func (s *Sender) Send(m *Message) error {
	return smtp.SendMail(fmt.Sprintf("%s:%s", host, port), s.auth, mail, m.To, m.ToBytes())
}

func NewMessage(s, b string) *Message {
	return &Message{Subject: s, Body: b, Attachment: []byte{}}
}

func (m *Message) AttachFile(src string) error {
	b, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	m.Attachment = b
	return nil
}

func (m *Message) ToBytes() []byte {
	buf := bytes.NewBuffer(nil)
	boundary := "skibidi-va-pa-pa"

	// Заголовки письма
	buf.WriteString(fmt.Sprintf("From: %s\r\n", mail))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(m.To, ", ")))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", m.Subject))
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", boundary))
	buf.WriteString("\r\n") // Пустая строка отделяет заголовки от тела

	// Текстовое сообщение
	buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	buf.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
	buf.WriteString("\r\n") // Пустая строка отделяет заголовки от тела
	buf.WriteString(m.Body)
	buf.WriteString("\r\n") // Заканчиваем тело сообщения

	// Фото
	if len(m.Attachment) > 0 {
		buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		buf.WriteString(fmt.Sprintf("Content-Type: %s\r\n", http.DetectContentType(m.Attachment)))
		buf.WriteString("Content-Transfer-Encoding: base64\r\n")
		buf.WriteString("Content-Disposition: attachment; filename=\"photo.jpg\"\r\n")
		buf.WriteString("\r\n") // Пустая строка отделяет заголовки от тела

		b := make([]byte, base64.StdEncoding.EncodedLen(len(m.Attachment)))
		base64.StdEncoding.Encode(b, m.Attachment)
		buf.Write(b)
		buf.WriteString("\r\n")
	}

	// Закрывающий boundary
	buf.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	return buf.Bytes()
}

func MailSender(fineData models.FineData) error {
	if fineData.Mail == "" {
		return fmt.Errorf("указана некорректная почта")
	}

	InitEmailConfig()

	sender := New()

	// Формирование сообщения о правонарушении
	m := NewMessage(
		"Уведомление о правонарушении",
		fmt.Sprintf("Вам назначается штраф в размере %d рублей.\n"+
			"Кординаты: %s\n"+
			"Тип и занчение правонарушения: %s, %s\n"+
			"Дата правонарушения: %s\n\n"+
			"Фото происшествия прилагаются:",
			fineData.Violation.Amount, fineData.Coordinated, fineData.Violation.Type, fineData.ViolationValue, fineData.Date,
		),
	)
	// Указание адрессанта
	m.To = []string{fineData.Mail}

	// Загрузка фото
	err := m.AttachFile(".." + fineData.PhotoUrl)
	if err != nil {
		return err
	}

	return sender.Send(m)
}
