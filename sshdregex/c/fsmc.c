
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

l1:
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c != 's') goto l0;

l2:
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c == 'h') goto l3;
	if (c == 's') goto l2;
	goto l0;

l3:
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c == 'd') goto l4;
	if (c == 's') goto l1;
	goto l0;

l4:
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c == '[') goto l5;
	if (c == 's') goto l1;
	goto l0;

l5:
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c <= '/') goto l0;
	if (c <= '9') goto l6;
	if (c <= 'r') goto l0;
	if (c == 's') goto l1;
	goto l0;

l6:
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c <= '/') goto l0;
	if (c <= '9') goto l7;
	if (c <= 'r') goto l0;
	if (c == 's') goto l1;
	goto l0;

l7:
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c <= '/') goto l0;
	if (c <= '9') goto l8;
	if (c <= 'r') goto l0;
	if (c == 's') goto l1;
	goto l0;

l8:
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c <= '/') goto l0;
	if (c <= '9') goto l9;
	if (c <= 'r') goto l0;
	if (c == 's') goto l1;
	goto l0;

l9:
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c <= '/') goto l0;
	if (c <= '9') goto l10;
	if (c <= 'r') goto l0;
	if (c == 's') goto l1;
	goto l0;

l10:
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c == ']') goto l11;
	if (c == 's') goto l1;
	goto l0;

l11:
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c == ':') goto l12;
	if (c == 's') goto l1;
	goto l0;

l12:
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c == '\t') goto l12;
	if (c == ' ') goto l12;
	if (c == 'F') goto l13;
	if (c == 's') goto l1;
	goto l0;

l13:
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c == 'a') goto l14;
	if (c == 's') goto l1;
	goto l0;

l14:
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c == 'i') goto l15;
	if (c == 's') goto l1;
	goto l0;

l15:
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c == 'l') goto l16;
	if (c == 's') goto l1;
	goto l0;

l16:
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c == 'e') goto l17;
	if (c == 's') goto l1;
	goto l0;

l17:
	if (p == e) return -1;

	c = (unsigned char) *p++;
	if (c == 'd') goto l18;
	if (c == 's') goto l1;
	goto l0;

l18:
	if (p == e) return 0x1; /* "sshd\[[0-9]{5}\]:[[:blank:]]*Failed" */

	c = (unsigned char) *p++;
	if (c != 's') goto l18;

l19:
	if (p == e) return 0x1; /* "sshd\[[0-9]{5}\]:[[:blank:]]*Failed" */

	c = (unsigned char) *p++;
	if (c != 's') goto l18;

l20:
	if (p == e) return 0x1; /* "sshd\[[0-9]{5}\]:[[:blank:]]*Failed" */

	c = (unsigned char) *p++;
	if (c == 'h') goto l21;
	if (c == 's') goto l20;
	goto l18;

l21:
	if (p == e) return 0x1; /* "sshd\[[0-9]{5}\]:[[:blank:]]*Failed" */

	c = (unsigned char) *p++;
	if (c == 'd') goto l22;
	if (c == 's') goto l19;
	goto l18;

l22:
	if (p == e) return 0x1; /* "sshd\[[0-9]{5}\]:[[:blank:]]*Failed" */

	c = (unsigned char) *p++;
	if (c == '[') goto l23;
	if (c == 's') goto l19;
	goto l18;

l23:
	if (p == e) return 0x1; /* "sshd\[[0-9]{5}\]:[[:blank:]]*Failed" */

	c = (unsigned char) *p++;
	if (c <= '/') goto l18;
	if (c <= '9') goto l24;
	if (c <= 'r') goto l18;
	if (c == 's') goto l19;
	goto l18;

l24:
	if (p == e) return 0x1; /* "sshd\[[0-9]{5}\]:[[:blank:]]*Failed" */

	c = (unsigned char) *p++;
	if (c <= '/') goto l18;
	if (c <= '9') goto l25;
	if (c <= 'r') goto l18;
	if (c == 's') goto l19;
	goto l18;

l25:
	if (p == e) return 0x1; /* "sshd\[[0-9]{5}\]:[[:blank:]]*Failed" */

	c = (unsigned char) *p++;
	if (c <= '/') goto l18;
	if (c <= '9') goto l26;
	if (c <= 'r') goto l18;
	if (c == 's') goto l19;
	goto l18;

l26:
	if (p == e) return 0x1; /* "sshd\[[0-9]{5}\]:[[:blank:]]*Failed" */

	c = (unsigned char) *p++;
	if (c <= '/') goto l18;
	if (c <= '9') goto l27;
	if (c <= 'r') goto l18;
	if (c == 's') goto l19;
	goto l18;

l27:
	if (p == e) return 0x1; /* "sshd\[[0-9]{5}\]:[[:blank:]]*Failed" */

	c = (unsigned char) *p++;
	if (c <= '/') goto l18;
	if (c <= '9') goto l28;
	if (c <= 'r') goto l18;
	if (c == 's') goto l19;
	goto l18;

l28:
	if (p == e) return 0x1; /* "sshd\[[0-9]{5}\]:[[:blank:]]*Failed" */

	c = (unsigned char) *p++;
	if (c == ']') goto l29;
	if (c == 's') goto l19;
	goto l18;

l29:
	if (p == e) return 0x1; /* "sshd\[[0-9]{5}\]:[[:blank:]]*Failed" */

	c = (unsigned char) *p++;
	if (c == ':') goto l30;
	if (c == 's') goto l19;
	goto l18;

l30:
	if (p == e) return 0x1; /* "sshd\[[0-9]{5}\]:[[:blank:]]*Failed" */

	c = (unsigned char) *p++;
	if (c == '\t') goto l30;
	if (c == ' ') goto l30;
	if (c == 'F') goto l31;
	if (c == 's') goto l19;
	goto l18;

l31:
	if (p == e) return 0x1; /* "sshd\[[0-9]{5}\]:[[:blank:]]*Failed" */

	c = (unsigned char) *p++;
	if (c == 'a') goto l32;
	if (c == 's') goto l19;
	goto l18;

l32:
	if (p == e) return 0x1; /* "sshd\[[0-9]{5}\]:[[:blank:]]*Failed" */

	c = (unsigned char) *p++;
	if (c == 'i') goto l33;
	if (c == 's') goto l19;
	goto l18;

l33:
	if (p == e) return 0x1; /* "sshd\[[0-9]{5}\]:[[:blank:]]*Failed" */

	c = (unsigned char) *p++;
	if (c == 'l') goto l34;
	if (c == 's') goto l19;
	goto l18;

l34:
	if (p == e) return 0x1; /* "sshd\[[0-9]{5}\]:[[:blank:]]*Failed" */

	c = (unsigned char) *p++;
	if (c == 's') goto l19;
	goto l18;
}

