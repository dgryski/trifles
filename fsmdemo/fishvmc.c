
int
fishmain(const char *b, const char *e)
{
	const char *p;
	int c;

	p = b;


l0:
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c != 'f') goto l0;

l1: /* e.g. "f" */
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c == 'f') goto l1;
	if (c != 'i') goto l0;

l2: /* e.g. "fi" */
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c == 'f') goto l1;
	if (c == 'i') goto l2;
	if (c != 's') goto l0;
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c == 'f') goto l1;
	if (c != 'h') goto l0;
	return 0x1; /* "fi+sh" */
}

