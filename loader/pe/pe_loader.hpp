#pragma once

#include <Windows.h>

#ifdef __cplusplus
extern "C" {
#endif
    void peLoader(unsigned char *data, const wchar_t* cmdline);
#ifdef __cplusplus
} 
#endif