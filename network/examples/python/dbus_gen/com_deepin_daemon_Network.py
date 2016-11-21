'''
Created with dbus2any

https://github.com/hugosenari/dbus2any


This code require python-dbus

Parameters:

* pydbusclient.tpl
* ./dbus_gen/dbus_dde_daemon_network.xml

See also:
    http://dbus.freedesktop.org/doc/dbus-specification.html
    http://dbus.freedesktop.org/doc/dbus-python/doc/tutorial.html
'''

import dbus



class Network(object):
    '''
    com.deepin.daemon.Network

    Usage:
    ------

    Instantiate this class and access the instance members and methods

    >>> obj = Network(BUS_NAME, OBJECT_PATH)

    '''

    def __init__(self, bus_name, object_path, interface=None, bus=None):
        '''Constructor'''
        self._dbus_interface_name = interface or "com.deepin.daemon.Network"
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

    
    def ActivateAccessPoint(self, *arg, **kw):
        '''
        Method (call me)
        params:
            : STRING
            : OBJECT_PATH
            : OBJECT_PATH
            
        return:
            : OBJECT_PATH
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.ActivateAccessPoint(*arg, **kw)

    def ActivateConnection(self, *arg, **kw):
        '''
        Method (call me)
        params:
            : STRING
            : OBJECT_PATH
            
        return:
            : OBJECT_PATH
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.ActivateConnection(*arg, **kw)

    def CancelSecret(self, *arg, **kw):
        '''
        Method (call me)
        params:
            : STRING
            : STRING
            
        
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.CancelSecret(*arg, **kw)

    def CreateConnection(self, *arg, **kw):
        '''
        Method (call me)
        params:
            : STRING
            : OBJECT_PATH
            
        return:
            : OBJECT_PATH
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.CreateConnection(*arg, **kw)

    def CreateConnectionForAccessPoint(self, *arg, **kw):
        '''
        Method (call me)
        params:
            : OBJECT_PATH
            : OBJECT_PATH
            
        return:
            : OBJECT_PATH
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.CreateConnectionForAccessPoint(*arg, **kw)

    def DeactivateConnection(self, *arg, **kw):
        '''
        Method (call me)
        params:
            : STRING
            
        
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.DeactivateConnection(*arg, **kw)

    def DeleteConnection(self, *arg, **kw):
        '''
        Method (call me)
        params:
            : STRING
            
        
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.DeleteConnection(*arg, **kw)

    def DisconnectDevice(self, *arg, **kw):
        '''
        Method (call me)
        params:
            : OBJECT_PATH
            
        
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.DisconnectDevice(*arg, **kw)

    def EditConnection(self, *arg, **kw):
        '''
        Method (call me)
        params:
            : STRING
            : OBJECT_PATH
            
        return:
            : OBJECT_PATH
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.EditConnection(*arg, **kw)

    def EnableDevice(self, *arg, **kw):
        '''
        Method (call me)
        params:
            : OBJECT_PATH
            : BOOLEAN
            
        
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.EnableDevice(*arg, **kw)

    def FeedSecret(self, *arg, **kw):
        '''
        Method (call me)
        params:
            : STRING
            : STRING
            : STRING
            : BOOLEAN
            
        
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.FeedSecret(*arg, **kw)

    def GetAccessPoints(self, *arg, **kw):
        '''
        Method (call me)
        params:
            : OBJECT_PATH
            
        return:
            : STRING
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.GetAccessPoints(*arg, **kw)

    def GetActiveConnectionInfo(self, *arg, **kw):
        '''
        Method (call me)
        return:
            : STRING
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.GetActiveConnectionInfo(*arg, **kw)

    def GetAutoProxy(self, *arg, **kw):
        '''
        Method (call me)
        return:
            : STRING
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.GetAutoProxy(*arg, **kw)

    def GetProxy(self, *arg, **kw):
        '''
        Method (call me)
        params:
            : STRING
            
        return:
            : STRING
            : STRING
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.GetProxy(*arg, **kw)

    def GetProxyIgnoreHosts(self, *arg, **kw):
        '''
        Method (call me)
        return:
            : STRING
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.GetProxyIgnoreHosts(*arg, **kw)

    def GetProxyMethod(self, *arg, **kw):
        '''
        Method (call me)
        return:
            : STRING
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.GetProxyMethod(*arg, **kw)

    def GetSupportedConnectionTypes(self, *arg, **kw):
        '''
        Method (call me)
        return:
            : as
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.GetSupportedConnectionTypes(*arg, **kw)

    def GetWiredConnectionUuid(self, *arg, **kw):
        '''
        Method (call me)
        params:
            : OBJECT_PATH
            
        return:
            : STRING
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.GetWiredConnectionUuid(*arg, **kw)

    def IsDeviceEnabled(self, *arg, **kw):
        '''
        Method (call me)
        params:
            : OBJECT_PATH
            
        return:
            : BOOLEAN
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.IsDeviceEnabled(*arg, **kw)

    def SetAutoProxy(self, *arg, **kw):
        '''
        Method (call me)
        params:
            : STRING
            
        
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.SetAutoProxy(*arg, **kw)

    def SetDeviceManaged(self, *arg, **kw):
        '''
        Method (call me)
        params:
            : STRING
            : BOOLEAN
            
        
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.SetDeviceManaged(*arg, **kw)

    def SetProxy(self, *arg, **kw):
        '''
        Method (call me)
        params:
            : STRING
            : STRING
            : STRING
            
        
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.SetProxy(*arg, **kw)

    def SetProxyIgnoreHosts(self, *arg, **kw):
        '''
        Method (call me)
        params:
            : STRING
            
        
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.SetProxyIgnoreHosts(*arg, **kw)

    def SetProxyMethod(self, *arg, **kw):
        '''
        Method (call me)
        params:
            : STRING
            
        
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._dbus_interface.SetProxyMethod(*arg, **kw)

    def NeedSecrets(self, callback):
        '''
        Signal (wait for me)
        callback params:
             STRING
             STRING
             STRING
             BOOLEAN
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448s
        '''
        self._dbus_interface.connect_to_signal('NeedSecrets', callback)
        return self

    def NeedSecretsFinished(self, callback):
        '''
        Signal (wait for me)
        callback params:
             STRING
             STRING
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448s
        '''
        self._dbus_interface.connect_to_signal('NeedSecretsFinished', callback)
        return self

    def AccessPointAdded(self, callback):
        '''
        Signal (wait for me)
        callback params:
             STRING
             STRING
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448s
        '''
        self._dbus_interface.connect_to_signal('AccessPointAdded', callback)
        return self

    def AccessPointRemoved(self, callback):
        '''
        Signal (wait for me)
        callback params:
             STRING
             STRING
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448s
        '''
        self._dbus_interface.connect_to_signal('AccessPointRemoved', callback)
        return self

    def AccessPointPropertiesChanged(self, callback):
        '''
        Signal (wait for me)
        callback params:
             STRING
             STRING
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448s
        '''
        self._dbus_interface.connect_to_signal('AccessPointPropertiesChanged', callback)
        return self

    def DeviceEnabled(self, callback):
        '''
        Signal (wait for me)
        callback params:
             STRING
             BOOLEAN
            
        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448s
        '''
        self._dbus_interface.connect_to_signal('DeviceEnabled', callback)
        return self

    @property
    def State(self):
        '''
        Property (acess me)
        Type:
            UINT32 read

        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._get_property('State')
    
    @property
    def NetworkingEnabled(self):
        '''
        Property (acess me)
        Type:
            BOOLEAN readwrite

        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._get_property('NetworkingEnabled')
    
    @NetworkingEnabled.setter
    def NetworkingEnabled(self, value):
        self._set_property('NetworkingEnabled', value)
    
    @property
    def VpnEnabled(self):
        '''
        Property (acess me)
        Type:
            BOOLEAN readwrite

        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._get_property('VpnEnabled')
    
    @VpnEnabled.setter
    def VpnEnabled(self, value):
        self._set_property('VpnEnabled', value)
    
    @property
    def Devices(self):
        '''
        Property (acess me)
        Type:
            STRING read

        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._get_property('Devices')
    
    @property
    def Connections(self):
        '''
        Property (acess me)
        Type:
            STRING read

        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._get_property('Connections')
    
    @property
    def ActiveConnections(self):
        '''
        Property (acess me)
        Type:
            STRING read

        See also:
            http://dbus.freedesktop.org/doc/dbus-specification.html#idp94392448
        '''
        return self._get_property('ActiveConnections')
    

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


