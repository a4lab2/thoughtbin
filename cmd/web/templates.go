package main

import "a4lab2.com/thoughtbin/pkg/models"

type templateData struct {
	Thought  *models.Thought
	Thoughts []*models.Thought
}
