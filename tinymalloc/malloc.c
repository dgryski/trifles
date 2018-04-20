/* the following is a malloc() implementation based on
 * http://www.concentric.net/~Ttwang/tech/smallmem.htm
 */

#define ALIGN(a,b)  (((a)+b-1)&~(b-1))

#define NUM_CHUNK_LISTS 5

#define BIG_ACCOUNTING_INFO 8

/* bits are numbered 31..0 */
#define ISSET(w,b)  ((w) & (1 << b))
#define SETBIT(w,b) w |= (1 << b)
#define CLRBIT(w,b) w &= ~(1 << b)

/* this fits on a page */
typedef struct __malloc_chunk_t {
   struct __malloc_chunk_t *next;
   unsigned short size;
   unsigned short bitmap;
   char data[120];
} __malloc_chunk_t;

typedef struct __malloc_big_t {
   int size;
   struct __malloc_big_t *next;
   struct __malloc_big_t *prev;
} __malloc_big_t;

static void *curheap = 0;

int
brk(void *end)
{
   return (int) Brk(end);
}

void *
sbrk(unsigned int ptrdiff)
{

   void *newend;

   /* have we been initialized yet? */
   if (curheap == 0) {
      curheap = (void *) brk(0);
   }

   /* are they just querying the current heap? */
   if (ptrdiff == 0) {
      return curheap;
   }

   /* where we want our new heap to end */
   newend = curheap + ptrdiff;

   /* can we get the memory? */
   if (brk(newend) < 0) {
      return 0;
   }

   /* we got the memory */
   curheap = newend;
   return curheap - ptrdiff;
}

int
get_chunk_list(int size)
{
   if (size <= 8)
      return 0;
   if (size <= 16)
      return 1;
   if (size <= 24)
      return 2;
   if (size <= 40)
      return 3;
   if (size <= 56)
      return 4;
   return -1;
}

int
get_chunk_size(int size)
{
   if (size <= 8)
      return 8;
   if (size <= 16)
      return 16;
   if (size <= 24)
      return 24;
   if (size <= 40)
      return 40;
   if (size <= 56)
      return 56;

   return -1;
}

static __malloc_chunk_t *malloc_chunks[NUM_CHUNK_LISTS];
static __malloc_chunk_t *free_chunks = 0;
static __malloc_big_t *big_head = 0;
static int malloc_initialized = 0;

static void
malloc_init(void)
{

   int i;

   for (i = 0; i < NUM_CHUNK_LISTS; i++)
      malloc_chunks[i] = 0;

   free_chunks = 0;
   big_head = 0;

   malloc_initialized = 1;
}

static void *
big_malloc(int size)
{

   __malloc_big_t *curr, *next;
   int sz;
   void *ptr;

   curr = big_head;

   /* the size of space we actually need */
   sz = ALIGN(size + BIG_ACCOUNTING_INFO, 4);

   /* search for a chunk we can use */
   while (curr && curr->size < sz) {
      curr = curr->next;
   }

   /* didn't find space -- add some space to the free list */
   if (!curr) {

      /* 522 = 4*128 (blocks) + 2*4 (sentinals) + 8*4 (accounting) */
      int sz1 = MAX(sz + 8, 552);

      curr = sbrk(sz1);
      /* couldn't get any more memory -- fail! */
      if (!curr)
         return 0;

      /* sentinals for free() merging */
      *(int *) curr = 0;
      *(int *) ((char *) curr + sz1 - 4) = 0;

      curr = (void *) ((char *) curr + 4);

      curr->size = sz1 - 8;
      curr->next = big_head;
      curr->prev = 0;

      if (big_head) {
         big_head->prev = curr;
      } else {
         big_head = curr;
      }
   }

   /* we can grab some memory from ``curr'' */

   /* is there enough leftover space to justify a new node in the free list? */
   if (curr->size <= (sz + sizeof(__malloc_big_t))) {
      /* nope -- allocate the whole thing */
      sz = curr->size;

      /* remove this block from the list */
      if (curr->next)
         curr->next->prev = curr->prev;

      if (curr->prev)
         curr->prev->next = curr->next;
      else
         big_head = curr->next; /* must be at head of list */

   } else {
      /* yup -- add a new node */
      next = (void *) ((char *) curr + sz);
      next->size = curr->size - sz;
      next->next = curr->next;
      next->prev = curr->prev;
      /* accounting info */
      *(int *) ((char *) next + next->size - 4) = next->size;
      /* add to list */
      if (next->next)
         next->next->prev = next;
      if (next->prev)
         next->prev->next = next;
      else
         big_head = next;       /* must be at head of list */
   }

   /* set the accounting info */
   *(int *) curr = -sz;
   *(int *) ((char *) curr + sz - 4) = -sz;
   ptr = (void *) ((char *) curr + 4);

   return ptr;
}

void
big_free(void *ptr)
{

   int *sz = (int *) ((char *) ptr - 4);
   int *prevsz = (int *) ((char *) ptr - 8);
   int *nextsz = (int *) ((char *) ptr - *sz - 4);
   __malloc_big_t *ch, *nch;

   if (*prevsz <= 0 && *nextsz <= 0) {

      /* can't merge with either block -- add ourselves to free list */

      /* create node */
      ch = (__malloc_big_t *) sz;
      ch->size = -*sz;
      /* accounting info for future merges */
      *(int *) ((char *) ch + ch->size - 4) = ch->size;

      /* add to head of list */
      ch->prev = 0;
      ch->next = big_head;
      if (big_head)
         big_head->prev = ch;
      big_head = ch;

      return;
   }

   ch = 0;
   if (*prevsz > 0) {
      /* block before us is free -- merge! */
      ch = (__malloc_big_t *) ((char *) prevsz - *prevsz + 4);  /* sz is still -ve */
      ch->size += -*sz;

      /* reset accounting info for end of new block */
      *(int *) ((char *) ch + ch->size - 4) = ch->size;
   };


   if (*nextsz > 0) {

      if (ch) {
         /* we merged with the block before us */
         /* and now want to merge with the block after us too */

         /* remove block after us from linked list */
         nch = (__malloc_big_t *) nextsz;

         if (nch->next)
            nch->next->prev = nch->prev;
         if (nch->prev)
            nch->prev->next = nch->next;
         else
            big_head = nch->next;       /* no prev? must be head of list */

         /* increase our size accordingly */
         ch->size += nch->size;
         /* reset accounting info for end of new block */
         *(int *) ((char *) ch + ch->size - 4) = ch->size;

         return;

      }

      /* merge with block after us */
      nch = (__malloc_big_t *) nextsz;
      /* create new node */
      ch = (__malloc_big_t *) sz;
      ch->size = -ch->size + nch->size;
      ch->next = nch->next;
      ch->prev = nch->prev;

      /* reset accounting info for end of new block */
      *(int *) ((char *) ch + ch->size - 4) = ch->size;

      /* update node ptrs */
      if (ch->next)
         ch->next->prev = ch;

      if (ch->prev)
         ch->prev->next = ch;
      else                      /* must have been head of list */
         big_head = ch;

   }

   return;
}


static void *
small_malloc(int size)
{

   __malloc_chunk_t *chunk;
   int i, validbits;
   unsigned int mask;

   size = get_chunk_size(size);
   i = get_chunk_list(size);
   chunk = malloc_chunks[i];

   /* how many bits are valid for this chunks bitmap? */
   validbits = (sizeof(__malloc_chunk_t) - 8) / size;
   mask = (1 << (validbits)) - 1;

   /* search through chunk list */
   for (;;) {

      if (chunk == 0) {

         /* out of chunks -- allocate a new one */

         if (free_chunks) {
            /* try and grab a recycled chunk */
            chunk = free_chunks;
            free_chunks = free_chunks->next;
            chunk->next = 0;
         } else {
            /* nothing on the free list */
            /* ask the system for some memory */
            chunk = sbrk(128);
            /* couldn't grab more memory -- must be out */
            if (chunk == 0)
               return 0;
         }

         /* no blocks currently allocated on this page */
         chunk->size = get_chunk_size(size);
         chunk->bitmap = 0;

         /* add to head of list */
         chunk->next = malloc_chunks[i];
         malloc_chunks[i] = chunk;
      }

      /* scan chunk's bitmap for a free block */
      if ((chunk->bitmap & mask) == mask) {
         /* all blocks are used -- move on to next one */
         chunk = chunk->next;
         continue;
      }

      /* at least one is free */
      for (i = 0; i < validbits; i++) {
         if (!ISSET(chunk->bitmap, i)) {
            /* found free block */
            SETBIT(chunk->bitmap, i);
            return (void *) &chunk->data[chunk->size * i];
         }
      }
   }
}


int
small_free(void *ptr)
{

   __malloc_chunk_t *chunk = 0;
   __malloc_chunk_t *prev, *curr;
   int i;

   /* see if the ptr is on one of our lists */
   for (i = 0; i < NUM_CHUNK_LISTS; i++) {

      chunk = malloc_chunks[i];
      while (chunk) {
         if ((void *) ptr > (void *) chunk && ptr < (void *) (chunk + 1)) {
            /* yup, it was in our memory region */
            goto found_it;
         }
         chunk = chunk->next;
      }
   }

   /* nope -- didn't find it */
   return 0;

 found_it:
   /* which bit should we clear? */
   i = ((int) ptr - (int) chunk - 8) / chunk->size;

   CLRBIT(chunk->bitmap, i);

   if (chunk->bitmap == 0) {

      /* we can add this chunk to the free list */

      i = get_chunk_list(chunk->size);
      prev = 0;
      curr = malloc_chunks[i];

      /* find the chunk in it's list */
      while (curr && curr != chunk) {
         prev = curr;
         curr = curr->next;
      }

      /* unlink */
      if (prev == 0) {
         /* chunk was at head of list */
         malloc_chunks[i] = curr->next;
      } else {
         prev->next = curr->next;
      }

      /* add to head of free list */
      chunk->next = free_chunks;
      free_chunks = chunk;
   }

   return 1;
}

void *
malloc(unsigned int size)
{

   if (!malloc_initialized)
      malloc_init();

   if (size == 0)
      return 0;                 /* so says ANSI */

   /* use the special allocator for small sizes */
   if (size <= 56)
      return small_malloc(size);

   /* use the standard free list for large sizes */
   return big_malloc(size);

}

void
free(void *ptr)
{

   /* let them free null */
   if (!ptr)
      return;

   /* see if it was a small malloc */
   if (small_free(ptr))
      return;

   /* must've been a big malloc */
   big_free(ptr);

}

void
__malloc_memstats(void)
{

   int i;

   __malloc_chunk_t *schunk;
   __malloc_big_t *bchunk;

   if (!malloc_initialized) {
      printf("malloc library not initialized\n");
      return;
   }

   printf("\n");

   /* dump the chunk lists */
   for (i = 0; i < NUM_CHUNK_LISTS; i++) {
      schunk = malloc_chunks[i];
      printf("chunk list:%d ", i);
      if (schunk) {
         printf("size: %d: ", schunk->size);
         while (schunk) {
            printf(" %x", schunk->bitmap);
            schunk = schunk->next;
         }
         printf("\n");
      } else {
         printf("empty\n");
      }
   }

   /* count the free list */
   i = 0;
   schunk = free_chunks;
   while (schunk) {
      i++;
      schunk = schunk->next;
   }

   printf("Small chunk free list: %d items\n", i);

   if (big_head) {
      printf("Big chunk free list:\n");
      bchunk = big_head;
      while (bchunk) {
         printf("\tAddr: %x\tsize: %d\n", bchunk, bchunk->size);
         bchunk = bchunk->next;
      }
   }
}
