package main

//go:nobounds
func MatchRagel(data []byte) bool {

	var scanner_start int = 0
	var _ = scanner_start
	var scanner_first_final int = 18
	var _ = scanner_first_final
	var scanner_error int = -1
	var _ = scanner_error
	var scanner_en_main int = 0
	var _ = scanner_en_main
	cs, p, pe, eof := 0, 0, len(data), len(data)

	_ = eof

	{
		cs = int(scanner_start)

	}
	{
		switch cs {
		case 0:
			goto st_case_0
		case 1:
			goto st_case_1
		case 2:
			goto st_case_2
		case 3:
			goto st_case_3
		case 4:
			goto st_case_4
		case 5:
			goto st_case_5
		case 6:
			goto st_case_6
		case 7:
			goto st_case_7
		case 8:
			goto st_case_8
		case 9:
			goto st_case_9
		case 10:
			goto st_case_10
		case 11:
			goto st_case_11
		case 12:
			goto st_case_12
		case 13:
			goto st_case_13
		case 14:
			goto st_case_14
		case 15:
			goto st_case_15
		case 16:
			goto st_case_16
		case 17:
			goto st_case_17
		case 18:
			goto st_case_18

		}
	_st0:
		p += 1
	st_case_0:
		if p == pe {
			goto _out0

		}
		if (data[p]) == 115 {
			goto _st1

		}
		goto _st0
	_st1:
		p += 1
	st_case_1:
		if p == pe {
			goto _out1

		}
		if (data[p]) == 115 {
			goto _st2

		}
		goto _st0
	_st2:
		p += 1
	st_case_2:
		if p == pe {
			goto _out2

		}
		switch data[p] {
		case 104:
			{
				goto _st3

			}
		case 115:
			{
				goto _st2

			}

		}
		goto _st0
	_st3:
		p += 1
	st_case_3:
		if p == pe {
			goto _out3

		}
		switch data[p] {
		case 100:
			{
				goto _st4

			}
		case 115:
			{
				goto _st1

			}

		}
		goto _st0
	_st4:
		p += 1
	st_case_4:
		if p == pe {
			goto _out4

		}
		switch data[p] {
		case 91:
			{
				goto _st5

			}
		case 115:
			{
				goto _st1

			}

		}
		goto _st0
	_st5:
		p += 1
	st_case_5:
		if p == pe {
			goto _out5

		}
		if (data[p]) == 115 {
			goto _st1

		}
		if 48 <= (data[p]) && (data[p]) <= 57 {
			goto _st6

		}
		goto _st0
	_st6:
		p += 1
	st_case_6:
		if p == pe {
			goto _out6

		}
		if (data[p]) == 115 {
			goto _st1

		}
		if 48 <= (data[p]) && (data[p]) <= 57 {
			goto _st7

		}
		goto _st0
	_st7:
		p += 1
	st_case_7:
		if p == pe {
			goto _out7

		}
		if (data[p]) == 115 {
			goto _st1

		}
		if 48 <= (data[p]) && (data[p]) <= 57 {
			goto _st8

		}
		goto _st0
	_st8:
		p += 1
	st_case_8:
		if p == pe {
			goto _out8

		}
		if (data[p]) == 115 {
			goto _st1

		}
		if 48 <= (data[p]) && (data[p]) <= 57 {
			goto _st9

		}
		goto _st0
	_st9:
		p += 1
	st_case_9:
		if p == pe {
			goto _out9

		}
		if (data[p]) == 115 {
			goto _st1

		}
		if 48 <= (data[p]) && (data[p]) <= 57 {
			goto _st10

		}
		goto _st0
	_st10:
		p += 1
	st_case_10:
		if p == pe {
			goto _out10

		}
		switch data[p] {
		case 93:
			{
				goto _st11

			}
		case 115:
			{
				goto _st1

			}

		}
		goto _st0
	_st11:
		p += 1
	st_case_11:
		if p == pe {
			goto _out11

		}
		switch data[p] {
		case 58:
			{
				goto _st12

			}
		case 115:
			{
				goto _st1

			}

		}
		goto _st0
	_st12:
		p += 1
	st_case_12:
		if p == pe {
			goto _out12

		}
		switch data[p] {
		case 32:
			{
				goto _st12

			}
		case 70:
			{
				goto _st13

			}
		case 115:
			{
				goto _st1

			}

		}
		if 9 <= (data[p]) && (data[p]) <= 13 {
			goto _st12

		}
		goto _st0
	_st13:
		p += 1
	st_case_13:
		if p == pe {
			goto _out13

		}
		switch data[p] {
		case 97:
			{
				goto _st14

			}
		case 115:
			{
				goto _st1

			}

		}
		goto _st0
	_st14:
		p += 1
	st_case_14:
		if p == pe {
			goto _out14

		}
		switch data[p] {
		case 105:
			{
				goto _st15

			}
		case 115:
			{
				goto _st1

			}

		}
		goto _st0
	_st15:
		p += 1
	st_case_15:
		if p == pe {
			goto _out15

		}
		switch data[p] {
		case 108:
			{
				goto _st16

			}
		case 115:
			{
				goto _st1

			}

		}
		goto _st0
	_st16:
		p += 1
	st_case_16:
		if p == pe {
			goto _out16

		}
		switch data[p] {
		case 101:
			{
				goto _st17

			}
		case 115:
			{
				goto _st1

			}

		}
		goto _st0
	_st17:
		p += 1
	st_case_17:
		if p == pe {
			goto _out17

		}
		switch data[p] {
		case 100:
			{
				goto _ctr18

			}
		case 115:
			{
				goto _st1

			}

		}
		goto _st0
	_ctr18:
		{
			return true
		}
		goto _st18
	_st18:
		p += 1
	st_case_18:
		if p == pe {
			goto _out18

		}
		if (data[p]) == 115 {
			goto _st1

		}
		goto _st0
	_out0:
		cs = 0
		goto _out
	_out1:
		cs = 1
		goto _out
	_out2:
		cs = 2
		goto _out
	_out3:
		cs = 3
		goto _out
	_out4:
		cs = 4
		goto _out
	_out5:
		cs = 5
		goto _out
	_out6:
		cs = 6
		goto _out
	_out7:
		cs = 7
		goto _out
	_out8:
		cs = 8
		goto _out
	_out9:
		cs = 9
		goto _out
	_out10:
		cs = 10
		goto _out
	_out11:
		cs = 11
		goto _out
	_out12:
		cs = 12
		goto _out
	_out13:
		cs = 13
		goto _out
	_out14:
		cs = 14
		goto _out
	_out15:
		cs = 15
		goto _out
	_out16:
		cs = 16
		goto _out
	_out17:
		cs = 17
		goto _out
	_out18:
		cs = 18
		goto _out
	_out:
		{

		}

	}
	return false
}
