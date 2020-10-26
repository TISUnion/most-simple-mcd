#include "utils.h"

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
  //导入当前路径
  PyRun_SimpleString("import sys");
  PyRun_SimpleString("sys.path.append('./')");
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

int getFuncArgsLen(PyObject *pyFunc)
{
  int isFunc = PyFunction_Check(pyFunc);
  if (!isFunc) {
    return 0;
  }
  PyObject *inspectModule = PyImport_ImportModule("inspect");
  PyObject *inspectDict = PyModule_GetDict(inspectModule);
  PyObject *getfullargspecFunc = PyDict_GetItemString(inspectDict, "getfullargspec");
  PyObject *getfullargspecFuncRes = PyObject_CallFunction(getfullargspecFunc, "O", pyFunc);
  PyObject *getfullargspecArgsList = PyObject_GetAttrString(getfullargspecFuncRes, "args");
  int listLen = PyList_Size(getfullargspecArgsList);
  Py_XDECREF(inspectModule);
  Py_XDECREF(inspectDict);
  Py_XDECREF(getfullargspecFunc);
  Py_XDECREF(getfullargspecFuncRes);
  Py_XDECREF(getfullargspecArgsList);
  return listLen;
}