// Example on how to test a end-to-end web app with Playwright/Golang
package main

import (
	"fmt"
	"log"
	"reflect"

	"github.com/playwright-community/playwright-go"
)

func assertErrorToNilf(msg string, err error) {
	if err != nil {
		log.Fatalf(msg, err)
	}
}

func assertEqual(expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		panic(fmt.Sprintf("%v does not equal %v", actual, expected))
	}
}

const todoName = "Test the code"

func main() {
	pw, err := playwright.Run()
	assertErrorToNilf("could not lauch playwright: %w", err)
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})

	assertErrorToNilf("could not launch Chromium: %w", err)
	context, err := browser.NewContext()
	assertErrorToNilf("could not create context: %w", err)
	page, err := context.NewPage()
	assertErrorToNilf("could not create page: %w", err)

	_, err = page.Goto("https://demo.playwright.dev/todomvc/")
	assertErrorToNilf("could not goto: %w", err)

	// Helper function to get the amount of todos on the page
	assertCountOfTodos := func(shouldBeCount int) {
		count, err := page.Locator("ul.todo-list > li").Count()
		assertErrorToNilf("could not determine todo list count: %w", err)
		assertEqual(shouldBeCount, count)
	}

	// Initially there should be 0 entries
	assertCountOfTodos(0)

	// Adding a todo entry
	// (click in the input, enter the todo title and press the enter key)
	assertErrorToNilf("could not click: %v", page.Locator("input.new-todo").Click())
	assertErrorToNilf("could not type: %v", page.Locator("input.new-todo").Fill(todoName))
	assertErrorToNilf("could not press: %v", page.Locator("input.new-todo").Press("Enter"))

	// After adding 1 there should be 1 entry in the list
	assertCountOfTodos(1)

	// Here we get the text in the first todo item to see if it's the same with the user
	// has entered
	firstEntryText, err := page.Locator("ul.todo-list > li:nth-child(1) label").Evaluate("el => el.textContent", nil)
	assertErrorToNilf("could not get first todo entry: %w", err)
	assertEqual(todoName, firstEntryText)

	// if we filter now for completed entries, there should be 1
	assertErrorToNilf("could not click: %v", page.GetByRole("link", playwright.PageGetByRoleOptions{
		Name: "Completed",
	}).Click())

	// Clear the list of completed entries, then it should be again 0
	assertErrorToNilf("could not click: %v", page.Locator("text=Clear completed").Click())
	assertCountOfTodos(0)

	assertErrorToNilf("could not close browser: %w", browser.Close())
	assertErrorToNilf("could not stop Playwright: %w", pw.Stop())
}
