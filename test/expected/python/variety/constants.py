#
# Autogenerated by Frugal Compiler (1.24.0)
#
# DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING
#

from thrift.Thrift import TType, TMessageType, TException, TApplicationException
from .ttypes import *

redef_const = 582
DEFAULT_ID = -1
other_default = -1
thirtyfour = 34
MAPCONSTANT = {
    "hello": "world",
    "goodnight": "moon",
}
ConstEvent1 = Event(**{
    "ID": -2,
    "Message": "first one",
})
ConstEvent2 = Event(**{
    "ID": -7,
    "Message": "second one",
})
NumsList = [
    2,
    4,
    7,
    1,
]
NumsSet = set([
    1,
    3,
    8,
    0,
])
MAPCONSTANT2 = {
    "hello": Event(**{
        "ID": -2,
        "Message": "first here",
    }),
}
bin_const = "hello"
true_constant = True
false_constant = False
const_hc = 2
evil_string = "thin'g\" \""
evil_string2 = "th'ing\"ad\"f"
const_lower = TestLowercase(**{
    "lowercaseInt": 2,
})
