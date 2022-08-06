#include "loader.h"
#include <stdio.h>
#include "pe_loader.hpp"
#include <Windows.h>
#include <stdlib.h>

void pe(unsigned char *image, const char* cmd) {
    peLoader(image, NULL);
}