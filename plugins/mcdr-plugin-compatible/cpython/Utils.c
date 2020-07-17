#include <Python.h>
#include "Utils.h"

PyObject *CreateClass(char *className, PyObject *classDict)
{
  if (classDict == null)
  {
    classDict = PyDict_New();
  }
  PyObject *pClassBases = PyTuple_New(0);
  Py_XDECREF(classDict);
  return PyObject_CallFunctionObjArgs((PyObject *)&PyType_Type, className, pClassBases, classDict, NULL);
}

bool PyVmStart()
{
  // 初始化Python虚拟机
  Py_Initialize();
  // 判断Python虚拟机是否启动成功
  if (Py_IsInitialized() == 0)
  {
    return false;
  }
  return true;
}

void PyVmEnd()
{
  // 退出Python虚拟机(须先成功启动python虚拟机)
  Py_Finalize();
}

bool hasAttr(PyObject *classInstance, char *attrName)
{
  int res = PyObject_HasAttrString(classInstance, attrName) if (res == 1)
  {
    return true;
  }
  return false;
}