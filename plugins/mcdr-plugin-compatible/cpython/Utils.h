#ifndef _UTILS_H_
#define _UTILS_H_
#include <Python.h>

// 新建一个class
PyObject *CreateClass(char *className, PyObject *classDict);

// 开启python虚拟机
bool PyVmStart();

// 关闭python虚拟机
void PyVmEnd();

// 判断对象是否存在属性
bool hasAttr(PyObject *classInstance, char *attrName)

#endif