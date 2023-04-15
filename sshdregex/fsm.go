package main

import _ "unsafe"

//go:nobounds
func Match(data []byte) int {
	var idx = ^uint(0)

l0:
	if idx++; idx >= uint(len(data)) {
		return -1
	}

	if data[idx] != 's' {
		goto l0
	}

l1: // e.g. "s"
	if idx++; idx >= uint(len(data)) {
		return -1
	}

	if data[idx] != 's' {
		goto l0
	}

l2: // e.g. "ss"
	if idx++; idx >= uint(len(data)) {
		return -1
	}

	if data[idx] == 'h' {
		goto l3
	}

	if data[idx] == 's' {
		goto l2
	}

	{
		goto l0
	}

l3: // e.g. "ssh"
	if idx++; idx >= uint(len(data)) {
		return -1
	}

	if data[idx] == 'd' {
		goto l4
	}

	if data[idx] == 's' {
		goto l1
	}

	{
		goto l0
	}

l4: // e.g. "sshd"
	if idx++; idx >= uint(len(data)) {
		return -1
	}

	if data[idx] == '[' {
		goto l5
	}

	if data[idx] == 's' {
		goto l1
	}

	{
		goto l0
	}

l5: // e.g. "sshd["
	if idx++; idx >= uint(len(data)) {
		return -1
	}

	if data[idx] <= '/' {
		goto l0
	}

	if data[idx] <= '9' {
		goto l6
	}

	if data[idx] <= 'r' {
		goto l0
	}

	if data[idx] == 's' {
		goto l1
	}

	{
		goto l0
	}

l6: // e.g. "sshd[0"
	if idx++; idx >= uint(len(data)) {
		return -1
	}

	if data[idx] <= '/' {
		goto l0
	}

	if data[idx] <= '9' {
		goto l7
	}

	if data[idx] <= 'r' {
		goto l0
	}

	if data[idx] == 's' {
		goto l1
	}

	{
		goto l0
	}

l7: // e.g. "sshd[00"
	if idx++; idx >= uint(len(data)) {
		return -1
	}

	if data[idx] <= '/' {
		goto l0
	}

	if data[idx] <= '9' {
		goto l8
	}

	if data[idx] <= 'r' {
		goto l0
	}

	if data[idx] == 's' {
		goto l1
	}

	{
		goto l0
	}

l8: // e.g. "sshd[000"
	if idx++; idx >= uint(len(data)) {
		return -1
	}

	if data[idx] <= '/' {
		goto l0
	}

	if data[idx] <= '9' {
		goto l9
	}

	if data[idx] <= 'r' {
		goto l0
	}

	if data[idx] == 's' {
		goto l1
	}

	{
		goto l0
	}

l9: // e.g. "sshd[0000"
	if idx++; idx >= uint(len(data)) {
		return -1
	}

	if data[idx] <= '/' {
		goto l0
	}

	if data[idx] <= '9' {
		goto l10
	}

	if data[idx] <= 'r' {
		goto l0
	}

	if data[idx] == 's' {
		goto l1
	}

	{
		goto l0
	}

l10: // e.g. "sshd[00000"
	if idx++; idx >= uint(len(data)) {
		return -1
	}

	if data[idx] == ']' {
		goto l11
	}

	if data[idx] == 's' {
		goto l1
	}

	{
		goto l0
	}

l11: // e.g. "sshd[00000]"
	if idx++; idx >= uint(len(data)) {
		return -1
	}

	if data[idx] == ':' {
		goto l12
	}

	if data[idx] == 's' {
		goto l1
	}

	{
		goto l0
	}

l12: // e.g. "sshd[00000]:"
	if idx++; idx >= uint(len(data)) {
		return -1
	}

	if data[idx] <= '\b' {
		goto l0
	}

	if data[idx] <= '\r' {
		goto l12
	}

	if data[idx] <= '\x1f' {
		goto l0
	}

	if data[idx] == ' ' {
		goto l12
	}

	if data[idx] <= 'E' {
		goto l0
	}

	if data[idx] == 'F' {
		goto l13
	}

	if data[idx] <= 'r' {
		goto l0
	}

	if data[idx] == 's' {
		goto l1
	}

	{
		goto l0
	}

l13: // e.g. "sshd[00000]:F"
	if idx++; idx >= uint(len(data)) {
		return -1
	}

	if data[idx] == 'a' {
		goto l14
	}

	if data[idx] == 's' {
		goto l1
	}

	{
		goto l0
	}

l14: // e.g. "sshd[00000]:Fa"
	if idx++; idx >= uint(len(data)) {
		return -1
	}

	if data[idx] == 'i' {
		goto l15
	}

	if data[idx] == 's' {
		goto l1
	}

	{
		goto l0
	}

l15: // e.g. "sshd[00000]:Fai"
	if idx++; idx >= uint(len(data)) {
		return -1
	}

	if data[idx] == 'l' {
		goto l16
	}

	if data[idx] == 's' {
		goto l1
	}

	{
		goto l0
	}

l16: // e.g. "sshd[00000]:Fail"
	if idx++; idx >= uint(len(data)) {
		return -1
	}

	if data[idx] == 'e' {
		goto l17
	}

	if data[idx] == 's' {
		goto l1
	}

	{
		goto l0
	}

l17: // e.g. "sshd[00000]:Faile"
	if idx++; idx >= uint(len(data)) {
		return -1
	}

	if data[idx] == 'd' {
		goto l18
	}

	if data[idx] == 's' {
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
