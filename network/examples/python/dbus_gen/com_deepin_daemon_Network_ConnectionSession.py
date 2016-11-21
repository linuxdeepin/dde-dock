'''
Created with dbus2any

https://github.com/hugosenari/dbus2any


This code require python-dbus

Parameters:

* pydbusclient.tpl
* ./dbus_gen/dbus_dde_daemon_network_connectionsession.xml

See also:
    http://dbus.freedesktop.org/doc/dbus-specification.html
    http://dbus.freedesktop.org/doc/dbus-python/doc/tutorial.html
'''

import dbus



class Properties(object):
    '''
    org.freedesktop.DBus.Properties

    Usage:
    ------

    Instantiate this class and access the instance members and methods

    >>> obj = Properties(BUS_NAME, OBJECT_PATH)

    '''

    def __init__(self, bus_name, object_path, interface=None, bus=None):
        '''Constructor'''
        self._dbus_interface_name = interface or "org.freedesktop.DBus.Properties"
        self._dbus_object_path = object_path 
        self._dbus_name = bus_name 

        bus = bus or dbus.SessionBus()
        self._dbus_object =  bus.get_object(self._dbus_name, self._dbus_object_path)
        self._dbus_interface = dbus.Interface(self._dbus_object,
            dbus_interface=self._dbus_interface_name)
        self._dbus_properties = obj = dbus.Interface(self._dbus_object,
            "org.freedesktop.DBus.Properties")

    def _get_property(self, name):
        return self._dbus_properties.Get(self._dbus_interface_name, name)

    def _set_property(self, name, val):
        return self._dbus_properties.Set(self._dbus_interface_name, name, val)

    
    def Get(self, *arg, **kw):
        '''
        Method (call me)
        params:
            : STRING
            : STRING
            
        return:
            : VARIANT
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.Get(*arg, **kw)

    def GetAll(self, *arg, **kw):
        '''
        Method (call me)
        params:
            : STRING
            
        return:
            : a{sv}
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.GetAll(*arg, **kw)

    def InterfaceName(self, *arg, **kw):
        '''
        Method (call me)
        return:
            : STRING
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.InterfaceName(*arg, **kw)

    def Set(self, *arg, **kw):
        '''
        Method (call me)
        params:
            : STRING
            : STRING
            : VARIANT
            
        
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.Set(*arg, **kw)

    def PropertiesChanged(self, callback):
        '''
        Signal (wait for me)
        callback params:
             STRING
             a{sv}
             as
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448s
        '''
        self._dbus_interface.connect_to_signal('PropertiesChanged', callback)
        return self


class LifeManager(object):
    '''
    com.deepin.DBus.LifeManager

    Usage:
    ------

    Instantiate this class and access the instance members and methods

    >>> obj = LifeManager(BUS_NAME, OBJECT_PATH)

    '''

    def __init__(self, bus_name, object_path, interface=None, bus=None):
        '''Constructor'''
        self._dbus_interface_name = interface or "com.deepin.DBus.LifeManager"
        self._dbus_object_path = object_path 
        self._dbus_name = bus_name 

        bus = bus or dbus.SessionBus()
        self._dbus_object =  bus.get_object(self._dbus_name, self._dbus_object_path)
        self._dbus_interface = dbus.Interface(self._dbus_object,
            dbus_interface=self._dbus_interface_name)
        self._dbus_properties = obj = dbus.Interface(self._dbus_object,
            "org.freedesktop.DBus.Properties")

    def _get_property(self, name):
        return self._dbus_properties.Get(self._dbus_interface_name, name)

    def _set_property(self, name, val):
        return self._dbus_properties.Set(self._dbus_interface_name, name, val)

    
    def InterfaceName(self, *arg, **kw):
        '''
        Method (call me)
        return:
            : STRING
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.InterfaceName(*arg, **kw)

    def Ref(self, *arg, **kw):
        '''
        Method (call me)
        
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.Ref(*arg, **kw)

    def Unref(self, *arg, **kw):
        '''
        Method (call me)
        
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.Unref(*arg, **kw)


class ConnectionSession(object):
    '''
    com.deepin.daemon.ConnectionSession

    Usage:
    ------

    Instantiate this class and access the instance members and methods

    >>> obj = ConnectionSession(BUS_NAME, OBJECT_PATH)

    '''

    def __init__(self, bus_name, object_path, interface=None, bus=None):
        '''Constructor'''
        self._dbus_interface_name = interface or "com.deepin.daemon.ConnectionSession"
        self._dbus_object_path = object_path 
        self._dbus_name = bus_name 

        bus = bus or dbus.SessionBus()
        self._dbus_object =  bus.get_object(self._dbus_name, self._dbus_object_path)
        self._dbus_interface = dbus.Interface(self._dbus_object,
            dbus_interface=self._dbus_interface_name)
        self._dbus_properties = obj = dbus.Interface(self._dbus_object,
            "org.freedesktop.DBus.Properties")

    def _get_property(self, name):
        return self._dbus_properties.Get(self._dbus_interface_name, name)

    def _set_property(self, name, val):
        return self._dbus_properties.Set(self._dbus_interface_name, name, val)

    
    def Close(self, *arg, **kw):
        '''
        Method (call me)
        
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.Close(*arg, **kw)

    def DebugGetConnectionData(self, *arg, **kw):
        '''
        Method (call me)
        return:
            : a{sa{sv}}
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.DebugGetConnectionData(*arg, **kw)

    def DebugGetErrors(self, *arg, **kw):
        '''
        Method (call me)
        return:
            : a{sa{ss}}
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.DebugGetErrors(*arg, **kw)

    def DebugListKeyDetail(self, *arg, **kw):
        '''
        Method (call me)
        return:
            : STRING
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.DebugListKeyDetail(*arg, **kw)

    def GetAllKeys(self, *arg, **kw):
        '''
        Method (call me)
        return:
            : STRING
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.GetAllKeys(*arg, **kw)

    def GetAvailableValues(self, *arg, **kw):
        '''
        Method (call me)
        params:
            : STRING
            : STRING
            
        return:
            : STRING
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.GetAvailableValues(*arg, **kw)

    def GetKey(self, *arg, **kw):
        '''
        Method (call me)
        params:
            : STRING
            : STRING
            
        return:
            : STRING
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.GetKey(*arg, **kw)

    def GetKeyName(self, *arg, **kw):
        '''
        Method (call me)
        params:
            : STRING
            : STRING
            
        return:
            : STRING
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.GetKeyName(*arg, **kw)

    def IsDefaultExpandedSection(self, *arg, **kw):
        '''
        Method (call me)
        params:
            : STRING
            
        return:
            : BOOLEAN
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.IsDefaultExpandedSection(*arg, **kw)

    def Save(self, *arg, **kw):
        '''
        Method (call me)
        return:
            : BOOLEAN
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.Save(*arg, **kw)

    def SetKey(self, *arg, **kw):
        '''
        Method (call me)
        params:
            : STRING
            : STRING
            : STRING
            
        
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.SetKey(*arg, **kw)

    def ConnectionDataChanged(self, callback):
        '''
        Signal (wait for me)
        callback params:
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448s
        '''
        self._dbus_interface.connect_to_signal('ConnectionDataChanged', callback)
        return self

    @property
    def ConnectionPath(self):
        '''
        Property (acess me)
        Type:
            OBJECT_PATH read

        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._get_property('ConnectionPath')
    
    @property
    def Uuid(self):
        '''
        Property (acess me)
        Type:
            STRING read

        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._get_property('Uuid')
    
    @property
    def Type(self):
        '''
        Property (acess me)
        Type:
            STRING read

        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._get_property('Type')
    
    @property
    def AllowDelete(self):
        '''
        Property (acess me)
        Type:
            BOOLEAN read

        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._get_property('AllowDelete')
    
    @property
    def AllowEditConnectionId(self):
        '''
        Property (acess me)
        Type:
            BOOLEAN read

        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._get_property('AllowEditConnectionId')
    
    @property
    def AvailableVirtualSections(self):
        '''
        Property (acess me)
        Type:
            as read

        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._get_property('AvailableVirtualSections')
    
    @property
    def AvailableSections(self):
        '''
        Property (acess me)
        Type:
            as read

        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._get_property('AvailableSections')
    
    @property
    def AvailableKeys(self):
        '''
        Property (acess me)
        Type:
            a{sas} read

        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._get_property('AvailableKeys')
    
    @property
    def Errors(self):
        '''
        Property (acess me)
        Type:
            a{sa{ss}} read

        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._get_property('Errors')
    

class Introspectable(object):
    '''
    org.freedesktop.DBus.Introspectable

    Usage:
    ------

    Instantiate this class and access the instance members and methods

    >>> obj = Introspectable(BUS_NAME, OBJECT_PATH)

    '''

    def __init__(self, bus_name, object_path, interface=None, bus=None):
        '''Constructor'''
        self._dbus_interface_name = interface or "org.freedesktop.DBus.Introspectable"
        self._dbus_object_path = object_path 
        self._dbus_name = bus_name 

        bus = bus or dbus.SessionBus()
        self._dbus_object =  bus.get_object(self._dbus_name, self._dbus_object_path)
        self._dbus_interface = dbus.Interface(self._dbus_object,
            dbus_interface=self._dbus_interface_name)
        self._dbus_properties = obj = dbus.Interface(self._dbus_object,
            "org.freedesktop.DBus.Properties")

    def _get_property(self, name):
        return self._dbus_properties.Get(self._dbus_interface_name, name)

    def _set_property(self, name, val):
        return self._dbus_properties.Set(self._dbus_interface_name, name, val)

    
    def InterfaceName(self, *arg, **kw):
        '''
        Method (call me)
        return:
            : STRING
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.InterfaceName(*arg, **kw)

    def Introspect(self, *arg, **kw):
        '''
        Method (call me)
        return:
            : STRING
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.Introspect(*arg, **kw)


