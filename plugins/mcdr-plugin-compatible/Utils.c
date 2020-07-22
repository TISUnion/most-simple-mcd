#include "Utils.h"

PyObject *CreateClass(char *className, PyObject *classDict)
{
    if (classDict == NULL)
    {
      classDict = PyDict_New();
    }
    PyObject *pClassBases = PyTuple_New(0);
    return PyObject_CallFunctionObjArgs((PyObject *)&PyType_Type, PyUnicode_FromString(className), pClassBases, classDict, NULL);
}

int PyVmStart()
{
  // 初始化Python虚拟机
  Py_Initialize();
  // 判断Python虚拟机是否启动成功
  if (Py_IsInitialized() == 0)
  {
    return 0;
  }
  return 1;
}

void PyVmEnd()
{
  // 退出Python虚拟机(须先成功启动python虚拟机)
  Py_Finalize();
}

int hasAttr(PyObject *classInstance, char *attrName)
{
  int res = PyObject_HasAttrString(classInstance, attrName);
  if (res == 1)
  {
    return 1;
  }
  return 0;
}

char *pyObj2string(PyObject *obj)
{
  PyObject *repr = PyObject_Repr(obj);
  PyObject *str = PyUnicode_AsEncodedString(repr, "utf-8", "~E~");
  char *bytes = PyBytes_AS_STRING(str);
  Py_XDECREF(repr);
  Py_XDECREF(str);
  return bytes;
}