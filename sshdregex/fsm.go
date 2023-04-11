package main

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

l1:
	if idx++; idx >= uint(len(data)) {
		return -1
	}

	if data[idx] != 's' {
		goto l0
	}

l2:
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

l3:
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

l4:
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

l5:
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

l6:
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

l7:
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

l8:
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

l9:
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

l10:
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

l11:
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

l12:
	if idx++; idx >= uint(len(data)) {
		return -1
	}

	if data[idx] == '\t' {
		goto l12
	}

	if data[idx] == ' ' {
		goto l12
	}

	if data[idx] == 'F' {
		goto l13
	}

	if data[idx] == 's' {
		goto l1
	}

	{
		goto l0
	}

l13:
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

l14:
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

l15:
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

l16:
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

l17:
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

l18:
	if idx++; idx >= uint(len(data)) {
		return 18
	}

	if data[idx] != 's' {
		goto l18
	}

l19:
	if idx++; idx >= uint(len(data)) {
		return 19
	}

	if data[idx] != 's' {
		goto l18
	}

l20:
	if idx++; idx >= uint(len(data)) {
		return 20
	}

	if data[idx] == 'h' {
		goto l21
	}

	if data[idx] == 's' {
		goto l20
	}

	{
		goto l18
	}

l21:
	if idx++; idx >= uint(len(data)) {
		return 21
	}

	if data[idx] == 'd' {
		goto l22
	}

	if data[idx] == 's' {
		goto l19
	}

	{
		goto l18
	}

l22:
	if idx++; idx >= uint(len(data)) {
		return 22
	}

	if data[idx] == '[' {
		goto l23
	}

	if data[idx] == 's' {
		goto l19
	}

	{
		goto l18
	}

l23:
	if idx++; idx >= uint(len(data)) {
		return 23
	}

	if data[idx] <= '/' {
		goto l18
	}

	if data[idx] <= '9' {
		goto l24
	}

	if data[idx] <= 'r' {
		goto l18
	}

	if data[idx] == 's' {
		goto l19
	}

	{
		goto l18
	}

l24:
	if idx++; idx >= uint(len(data)) {
		return 24
	}

	if data[idx] <= '/' {
		goto l18
	}

	if data[idx] <= '9' {
		goto l25
	}

	if data[idx] <= 'r' {
		goto l18
	}

	if data[idx] == 's' {
		goto l19
	}

	{
		goto l18
	}

l25:
	if idx++; idx >= uint(len(data)) {
		return 25
	}

	if data[idx] <= '/' {
		goto l18
	}

	if data[idx] <= '9' {
		goto l26
	}

	if data[idx] <= 'r' {
		goto l18
	}

	if data[idx] == 's' {
		goto l19
	}

	{
		goto l18
	}

l26:
	if idx++; idx >= uint(len(data)) {
		return 26
	}

	if data[idx] <= '/' {
		goto l18
	}

	if data[idx] <= '9' {
		goto l27
	}

	if data[idx] <= 'r' {
		goto l18
	}

	if data[idx] == 's' {
		goto l19
	}

	{
		goto l18
	}

l27:
	if idx++; idx >= uint(len(data)) {
		return 27
	}

	if data[idx] <= '/' {
		goto l18
	}

	if data[idx] <= '9' {
		goto l28
	}

	if data[idx] <= 'r' {
		goto l18
	}

	if data[idx] == 's' {
		goto l19
	}

	{
		goto l18
	}

l28:
	if idx++; idx >= uint(len(data)) {
		return 28
	}

	if data[idx] == ']' {
		goto l29
	}

	if data[idx] == 's' {
		goto l19
	}

	{
		goto l18
	}

l29:
	if idx++; idx >= uint(len(data)) {
		return 29
	}

	if data[idx] == ':' {
		goto l30
	}

	if data[idx] == 's' {
		goto l19
	}

	{
		goto l18
	}

l30:
	if idx++; idx >= uint(len(data)) {
		return 30
	}

	if data[idx] == '\t' {
		goto l30
	}

	if data[idx] == ' ' {
		goto l30
	}

	if data[idx] == 'F' {
		goto l31
	}

	if data[idx] == 's' {
		goto l19
	}

	{
		goto l18
	}

l31:
	if idx++; idx >= uint(len(data)) {
		return 31
	}

	if data[idx] == 'a' {
		goto l32
	}

	if data[idx] == 's' {
		goto l19
	}

	{
		goto l18
	}

l32:
	if idx++; idx >= uint(len(data)) {
		return 32
	}

	if data[idx] == 'i' {
		goto l33
	}

	if data[idx] == 's' {
		goto l19
	}

	{
		goto l18
	}

l33:
	if idx++; idx >= uint(len(data)) {
		return 33
	}

	if data[idx] == 'l' {
		goto l34
	}

	if data[idx] == 's' {
		goto l19
	}

	{
		goto l18
	}

l34:
	if idx++; idx >= uint(len(data)) {
		return 34
	}

	if data[idx] == 's' {
		goto l19
	}

	{
		goto l18
	}

}
