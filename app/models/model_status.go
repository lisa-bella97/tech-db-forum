/*
 * forum
 *
 * Тестовое задание для реализации проекта \"Форумы\" на курсе по базам данных в Технопарке Mail.ru (https://park.mail.ru).
 *
 * API version: 0.1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package models

type Status struct {

	// Кол-во пользователей в базе данных.
	User int32 `json:"user"`

	// Кол-во разделов в базе данных.
	Forum int32 `json:"forum"`

	// Кол-во веток обсуждения в базе данных.
	Thread int32 `json:"thread"`

	// Кол-во сообщений в базе данных.
	Post int64 `json:"post"`
}
