#include <stdio.h>
#include <stdlib.h>
#include <dlfcn.h>

int
main(int argc, char **argv)
{
    void *handle;
    int (*match)(char *s, int len);
    char *error;

   handle = dlopen("fsmso.so", RTLD_LAZY);
    if (!handle) {
        fprintf(stderr, "%s\n", dlerror());
        exit(EXIT_FAILURE);
    }

   dlerror();    /* Clear any existing error */

   *(void **) (&match) = dlsym(handle, "fishasmMatch");

   if ((error = dlerror()) != NULL)  {
        fprintf(stderr, "%s\n", error);
        exit(EXIT_FAILURE);
    }

    printf("%d\n", match("catfishies", 10));
    printf("%d\n", match("dogo", 4));
    dlclose(handle);
    exit(EXIT_SUCCESS);
}
