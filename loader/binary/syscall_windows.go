package main

import "syscall"

func showWindow(show bool) {
	if show {
		return
	}
	getWin := syscall.NewLazyDLL(string([]byte{'k', 'e', 'r', 'n', 'e', 'l', '3', '2'})).NewProc("GetConsoleWindow")
	showWin := syscall.NewLazyDLL(string([]byte{'u', 's', 'e', 'r', '3', '2'})).NewProc("ShowWindow")
	hwnd, _, _ := getWin.Call()
	if getWin == nil {
		return
	}

	var SW_HIDE uintptr = 0
	showWin.Call(hwnd, SW_HIDE)
}
