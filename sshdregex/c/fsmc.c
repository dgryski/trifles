
int
sshdmain(const char *b, const char *e)
{
	const char *p;
	int c;

	p = b;


l0:
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c != 's') goto l0;

l1: /* e.g. "s" */
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c != 's') goto l0;

l2: /* e.g. "ss" */
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c == 'h') goto l3;
	if (c == 's') goto l2;
	goto l0;

l3: /* e.g. "ssh" */
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c == 'd') goto l4;
	if (c == 's') goto l1;
	goto l0;

l4: /* e.g. "sshd" */
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c == '[') goto l5;
	if (c == 's') goto l1;
	goto l0;

l5: /* e.g. "sshd[" */
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c <= '/') goto l0;
	if (c <= '9') goto l6;
	if (c <= 'r') goto l0;
	if (c == 's') goto l1;
	goto l0;

l6: /* e.g. "sshd[0" */
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c <= '/') goto l0;
	if (c <= '9') goto l7;
	if (c <= 'r') goto l0;
	if (c == 's') goto l1;
	goto l0;

l7: /* e.g. "sshd[00" */
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c <= '/') goto l0;
	if (c <= '9') goto l8;
	if (c <= 'r') goto l0;
	if (c == 's') goto l1;
	goto l0;

l8: /* e.g. "sshd[000" */
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c <= '/') goto l0;
	if (c <= '9') goto l9;
	if (c <= 'r') goto l0;
	if (c == 's') goto l1;
	goto l0;

l9: /* e.g. "sshd[0000" */
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c <= '/') goto l0;
	if (c <= '9') goto l10;
	if (c <= 'r') goto l0;
	if (c == 's') goto l1;
	goto l0;

l10: /* e.g. "sshd[00000" */
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c == ']') goto l11;
	if (c == 's') goto l1;
	goto l0;

l11: /* e.g. "sshd[00000]" */
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c == ':') goto l12;
	if (c == 's') goto l1;
	goto l0;

l12: /* e.g. "sshd[00000]:" */
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c <= '\b') goto l0;
	if (c <= '\r') goto l12;
	if (c <= '\x1f') goto l0;
	if (c == ' ') goto l12;
	if (c <= 'E') goto l0;
	if (c == 'F') goto l13;
	if (c <= 'r') goto l0;
	if (c == 's') goto l1;
	goto l0;

l13: /* e.g. "sshd[00000]:F" */
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c == 'a') goto l14;
	if (c == 's') goto l1;
	goto l0;

l14: /* e.g. "sshd[00000]:Fa" */
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c == 'i') goto l15;
	if (c == 's') goto l1;
	goto l0;

l15: /* e.g. "sshd[00000]:Fai" */
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c == 'l') goto l16;
	if (c == 's') goto l1;
	goto l0;

l16: /* e.g. "sshd[00000]:Fail" */
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c == 'e') goto l17;
	if (c == 's') goto l1;
	goto l0;

l17: /* e.g. "sshd[00000]:Faile" */
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c == 'd') goto l18;
	if (c == 's') goto l1;
	goto l0;

l18: /* e.g. "sshd[00000]:Failed" */
	return 0x1; /* "sshd\[\d{5}\]:\s*Failed" */
}

