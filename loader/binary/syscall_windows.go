package main

import "syscall"

func hide(show bool) {
	if show {
		return
	}
	getWin := syscall.NewLazyDLL(string([]byte{'k', 'e', 'r', 'n', 'e', 'l', '3', '2'})).NewProc("GetConsoleWindow")
	showWin := syscall.NewLazyDLL(string([]byte{'u', 's', 'e', 'r', '3', '2'})).NewProc("ShowWindow")
	hwnd, _, _ := getWin.Call()
	if getWin == nil {
		return
	}
	if show {
		var SW_RESTORE uintptr = 9
		showWin.Call(hwnd, SW_RESTORE)
	} else {
		var SW_HIDE uintptr = 0
		showWin.Call(hwnd, SW_HIDE)
	}
}
