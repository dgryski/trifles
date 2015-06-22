package main

/*
#cgo pkg-config: luajit

#include <lua.h>
#include <lualib.h>
#include <lauxlib.h>
#include <stdlib.h>

static int luaExecute(lua_State *L) { return lua_pcall(L, 0, LUA_MULTRET, 0); }
static void luaNewtable(lua_State *L) { lua_createtable(L, 0, 0); }
static void luaSetglobal(lua_State *L,char *s)  { lua_setfield(L, LUA_GLOBALSINDEX, (s)); }

int luaopen_cmsgpack (lua_State *L);

*/
import "C"

import (
	"fmt"
	"unsafe"

	"gopkg.in/vmihailenco/msgpack.v2"
)

var sum = `
sum = function(t)
    total = 0
    for i, v in ipairs(t) do total = total + v end
    return total
end


data = cmsgpack.unpack(mpdata)
return sum(data)
`

func cptr(a []byte) *C.char { return (*C.char)(unsafe.Pointer(&a[0])) }

func main() {

	L := C.luaL_newstate()
	defer C.lua_close(L)

	C.luaL_openlibs(L) /* Load Lua libraries */
	C.luaopen_cmsgpack(L)

	sumS := C.CString(sum)
	defer C.free(unsafe.Pointer(sumS))

	C.luaL_loadstring(L, sumS)

	msg, _ := msgpack.Marshal([]int{1, 2, 3, 4, 5, 6})
	C.lua_pushstring(L, cptr(msg))

	dstr := C.CString("mpdata")
	defer C.free(unsafe.Pointer(dstr))
	C.luaSetglobal(L, dstr)

	C.luaExecute(L)

	sum := C.lua_tonumber(L, -1)

	fmt.Println(sum) // Output: 21
}
