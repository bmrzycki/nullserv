#include <stdio.h>

int main(int argc, char *argv[]) {
     FILE* f;
     size_t f_sz;
     int i, c;

     if (argc < 2) {
	  fprintf(stderr, "usage: file2gobyte file\n");
	  return 1;
     }

     f = fopen(argv[1], "rb");
     if (f == NULL) {
	  fprintf(stderr, "cannot open file %s\n", argv[1]);
	  return 1;
     }

     fseek(f, 0, SEEK_END);
     f_sz = ftell(f);
     rewind(f);

     printf("fileName := \"%s\"\n", argv[1]);
     printf("fileSize := %ld\n", f_sz);
     printf("fileData := []byte{\n\t\t");

     for (i=0; i < f_sz; i++) {
	  c = fgetc(f);
	  printf("'\\x%02x',", c);
	  if ((i + 1) % 8 == 0)
	       printf("\n\t\t");
	  else
	       printf(" ");
     }

     fclose(f);
     printf("\n}\n");
     return 0;
}
