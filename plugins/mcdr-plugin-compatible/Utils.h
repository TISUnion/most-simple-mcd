#ifndef _UTILS_H_
#define _UTILS_H_ 1
#include <Python.h>

typedef struct Server{
	char *name;
    char *id;
    int memory;
    int port;
    char *version;
    char *side;
    char *comment;
} Server;

typedef struct Info{
	int id;
	int hour;
	int min;
	int sec;
	char *raw_content;
	char *content;
	char *player;
	int source;
	char *logging_level;
	int is_player;
	int is_user;
} Info;

// 事件枚举
enum
{
      OnLoad = 1, OnUnload, OnInfo, OnUserInfo, OnPlayerJoined
};

// 新建一个class
PyObject *CreateClass(char *className, PyObject *classDict);

// 开启python虚拟机
int PyVmStart();

// 关闭python虚拟机
void PyVmEnd();

// 判断对象是否存在属性
int hasAttr(PyObject *classInstance, char *attrName);

// python对象转char*
char *pyObj2string(PyObject *obj);

// 获取python函数形参个数
int getFuncArgsLen(PyObject *pyFunc)
#endif