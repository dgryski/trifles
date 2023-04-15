package main

import "unsafe"

//go:nobounds
func UnsafeMatch(data []byte) int {
	var ptr, end *byte
	if len(data) > 0 {
		ptr = (*byte)(unsafe.Pointer(&data[0]))
		end = (*byte)(unsafe.Add(unsafe.Pointer(ptr), len(data)))
	}
	var c byte

l0:
	if ptr == end {
		return -1
	}
	ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
	c = *ptr

	if c != 's' {
		goto l0
	}

l1: // e.g. "s"
	if ptr == end {
		return -1
	}
	ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
	c = *ptr

	if c != 's' {
		goto l0
	}

l2: // e.g. "ss"
	if ptr == end {
		return -1
	}
	ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
	c = *ptr

	if c == 'h' {
		goto l3
	}

	if c == 's' {
		goto l2
	}

	{
		goto l0
	}

l3: // e.g. "ssh"
	if ptr == end {
		return -1
	}
	ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
	c = *ptr

	if c == 'd' {
		goto l4
	}

	if c == 's' {
		goto l1
	}

	{
		goto l0
	}

l4: // e.g. "sshd"
	if ptr == end {
		return -1
	}
	ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
	c = *ptr

	if c == '[' {
		goto l5
	}

	if c == 's' {
		goto l1
	}

	{
		goto l0
	}

l5: // e.g. "sshd["
	if ptr == end {
		return -1
	}
	ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
	c = *ptr

	if c <= '/' {
		goto l0
	}

	if c <= '9' {
		goto l6
	}

	if c <= 'r' {
		goto l0
	}

	if c == 's' {
		goto l1
	}

	{
		goto l0
	}

l6: // e.g. "sshd[0"
	if ptr == end {
		return -1
	}
	ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
	c = *ptr

	if c <= '/' {
		goto l0
	}

	if c <= '9' {
		goto l7
	}

	if c <= 'r' {
		goto l0
	}

	if c == 's' {
		goto l1
	}

	{
		goto l0
	}

l7: // e.g. "sshd[00"
	if ptr == end {
		return -1
	}
	ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
	c = *ptr

	if c <= '/' {
		goto l0
	}

	if c <= '9' {
		goto l8
	}

	if c <= 'r' {
		goto l0
	}

	if c == 's' {
		goto l1
	}

	{
		goto l0
	}

l8: // e.g. "sshd[000"
	if ptr == end {
		return -1
	}
	ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
	c = *ptr

	if c <= '/' {
		goto l0
	}

	if c <= '9' {
		goto l9
	}

	if c <= 'r' {
		goto l0
	}

	if c == 's' {
		goto l1
	}

	{
		goto l0
	}

l9: // e.g. "sshd[0000"
	if ptr == end {
		return -1
	}
	ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
	c = *ptr

	if c <= '/' {
		goto l0
	}

	if c <= '9' {
		goto l10
	}

	if c <= 'r' {
		goto l0
	}

	if c == 's' {
		goto l1
	}

	{
		goto l0
	}

l10: // e.g. "sshd[00000"
	if ptr == end {
		return -1
	}
	ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
	c = *ptr

	if c == ']' {
		goto l11
	}

	if c == 's' {
		goto l1
	}

	{
		goto l0
	}

l11: // e.g. "sshd[00000]"
	if ptr == end {
		return -1
	}
	ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
	c = *ptr

	if c == ':' {
		goto l12
	}

	if c == 's' {
		goto l1
	}

	{
		goto l0
	}

l12: // e.g. "sshd[00000]:"
	if ptr == end {
		return -1
	}
	ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
	c = *ptr

	if c <= '\b' {
		goto l0
	}

	if c <= '\r' {
		goto l12
	}

	if c <= '\x1f' {
		goto l0
	}

	if c == ' ' {
		goto l12
	}

	if c <= 'E' {
		goto l0
	}

	if c == 'F' {
		goto l13
	}

	if c <= 'r' {
		goto l0
	}

	if c == 's' {
		goto l1
	}

	{
		goto l0
	}

l13: // e.g. "sshd[00000]:F"
	if ptr == end {
		return -1
	}
	ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
	c = *ptr

	if c == 'a' {
		goto l14
	}

	if c == 's' {
		goto l1
	}

	{
		goto l0
	}

l14: // e.g. "sshd[00000]:Fa"
	if ptr == end {
		return -1
	}
	ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
	c = *ptr

	if c == 'i' {
		goto l15
	}

	if c == 's' {
		goto l1
	}

	{
		goto l0
	}

l15: // e.g. "sshd[00000]:Fai"
	if ptr == end {
		return -1
	}
	ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
	c = *ptr

	if c == 'l' {
		goto l16
	}

	if c == 's' {
		goto l1
	}

	{
		goto l0
	}

l16: // e.g. "sshd[00000]:Fail"
	if ptr == end {
		return -1
	}
	ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
	c = *ptr

	if c == 'e' {
		goto l17
	}

	if c == 's' {
		goto l1
	}

	{
		goto l0
	}

l17: // e.g. "sshd[00000]:Faile"
	if ptr == end {
		return -1
	}
	ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
	c = *ptr

	if c == 'd' {
		goto l18
	}

	if c == 's' {
		goto l1
	}

	{
		goto l0
	}

l18: // e.g. "sshd[00000]:Failed"
	{
		return 18
	}

}
