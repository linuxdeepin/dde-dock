#!/usr/bin/env python3

# Copyright (C) 2016 Deepin Technology Co., Ltd.
#
# This program is free software; you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation; either version 3 of the License, or
# (at your option) any later version.

# See also: ${networkmanager-src}/libnm/generate-setting-docs.py

import gi
gi.require_version('NM', '1.0')
from gi.repository import NM, GObject
import argparse, datetime, re, sys
import xml.etree.ElementTree as ET

dbus_type_name_map = {
    'b': 'ktypeBoolean', # boolean
    's': 'ktypeString',  # string
    'i': 'ktypeInt32',   # int32
    'u': 'ktypeUint32',  # uint32
    't': 'ktypeUint64',  # uint64
    'x': 'ktypeInt64',   # int64
    'y': 'ktypeByte',    # byte
    'as': 'ktypeArrayString',         # array of string
    'au': 'ktypeArrayUint32',         # array of uint32
    'ay': 'ktypeArrayByte',           # byte array
    'a{ss}': 'ktypeDictStringString', # dict of string to string
    'a{sv}': 'ktypeUnknown',          # vardict
    'aa{sv}': 'ktypeUnknown',          # vardict
    'aau': 'ktypeArrayArrayUint32',   # array of array of uint32
    'aay': 'ktypeArrayArrayByte',     # array of byte array
    'a(ayuay)': 'ktypeIpv6Addresses', # array of legacy IPv6 address struct
    'a(ayuayu)': 'ktypeIpv6Routes',   # array of legacy IPv6 route struct
}

ns_map = {
    'c':    'http://www.gtk.org/introspection/c/1.0',
    'gi':   'http://www.gtk.org/introspection/core/1.0',
    'glib': 'http://www.gtk.org/introspection/glib/1.0'
}

def get_prop_type(setting, pspec, propxml):
    dbus_type = setting.get_dbus_property_type(pspec.name).dup_string()
    prop_type = dbus_type_name_map[dbus_type]
    return prop_type

def get_default_value(setting, pspec, propxml):
    default_value = setting.get_property(pspec.name.replace('-', '_'))
    value_type = get_prop_type(setting, pspec, propxml)

    # fix default value
    if default_value is not None:
        if value_type == 'ktypeBoolean':
            default_value = str(default_value).lower()
        elif value_type == 'ktypeString':
            default_value = "'%s'" % default_value
        elif value_type == 'ktypeByte':
            default_value = "'%s'" % default_value

        if str(default_value).startswith("<") or str(default_value).startswith("'<"):
            default_value = None

    # set a default value if not exists
    if default_value is None:
        if value_type == 'ktypeBoolean':
            default_value = "False"
        elif value_type == 'ktypeString':
            default_value = "''"
        elif value_type == 'ktypeInt32':
            default_value = "0"
        elif value_type == 'ktypeUint32':
            default_value = "0"
        elif value_type == 'ktypeUint64':
            default_value = "0"
        elif value_type == 'ktypeInt64':
            default_value = "0"
        elif value_type == 'ktypeByte':
            default_value = "0"
        elif value_type == 'ktypeArrayString':
            default_value = '[]'
        elif value_type == 'ktypeArrayUint32':
            default_value = '[]'
        elif value_type == 'ktypeArrayByte':
            default_value = '[]'
        elif value_type == 'ktypeDictStringString':
            default_value = '{}'
        elif value_type == 'ktypeArrayArrayUint32':
            default_value = '[]'
        elif value_type == 'ktypeArrayArrayByte':
            default_value = '[]'
        elif value_type == 'ktypeIpv6Addresses':
            default_value = '[]'
        elif value_type == 'ktypeIpv6Routes':
            default_value = '[]'

    # wrap all value as string
    if default_value is not None:
        default_value = '"%s"' % default_value

    return default_value

def usage():
    print("Usage: %s --gir FILE --output FILE" % sys.argv[0])
    exit()

# Main loop
parser = argparse.ArgumentParser()
parser.add_argument('-g', '--gir', metavar='FILE', help='NM-1.0.gir file (default: /usr/share/gir-1.0/NM-1.0.gir)', default='/usr/share/gir-1.0/NM-1.0.gir')
parser.add_argument('-o', '--output', metavar='FILE', help='output file (default: ./nm_consts_gen.yml)', default='./nm_consts_gen.yml')

args = parser.parse_args()
if args.gir is None or args.output is None:
    usage()

girxml = ET.parse(args.gir).getroot()
outfile = open(args.output, mode='w')

basexml = girxml.find('./gi:namespace/gi:class[@name="Setting"]', ns_map)
settings = girxml.findall('./gi:namespace/gi:class[@parent="Setting"]', ns_map)
# HACK: Need a better way to do this
ipxml = girxml.find('./gi:namespace/gi:class[@name="SettingIPConfig"]', ns_map)
settings.extend(girxml.findall('./gi:namespace/gi:class[@parent="SettingIPConfig"]', ns_map))
settings = sorted(settings, key=lambda setting: setting.attrib['{%s}symbol-prefix' % ns_map['c']])

# generate setting keys
constants = {}
outfile.write("---\n")
outfile.write("NMSettings:\n")
for settingxml in settings:
    if 'abstract' in settingxml.attrib:
        continue

    new_func = NM.__getattr__(settingxml.attrib['name'])
    setting = new_func()

    setting_capcase_name = settingxml.attrib['name']
    setting_name_prefix = "NM_" + settingxml.attrib['{%s}symbol-prefix' % ns_map['c']].upper().replace('-', '_')
    outfile.write("  - SettingClass: %s\n" % settingxml.attrib['name'])
    setting_name = "%s_SETTING_NAME" % setting_name_prefix
    outfile.write("    Name: %s\n" % setting_name)
    outfile.write("    Value: %s\n" % setting.props.name)
    outfile.write("    Keys:\n")
    constants[setting_name]=setting.props.name

    setting_properties = { prop.name: prop for prop in GObject.list_properties(setting) }
    properties = sorted(set(setting_properties.keys()))

    for prop in properties:
        if prop == 'name':
            continue

        value_type = None
        default_value = None

        if prop in setting_properties:
            pspec = setting_properties[prop]
            propxml = settingxml.find('./gi:property[@name="%s"]' % pspec.name, ns_map)
            if propxml is None:
                propxml = basexml.find('./gi:property[@name="%s"]' % pspec.name, ns_map)
            if propxml is None:
                propxml = ipxml.find('./gi:property[@name="%s"]' % pspec.name, ns_map)

            value_type = get_prop_type(setting, pspec, propxml)
            default_value = get_default_value(setting, pspec, propxml)

        # ignore unknown type
        if value_type == "ktypeUnknown":
            continue

        key_name = "%s_%s" % (setting_name_prefix, prop.upper().replace('-', '_'))
        outfile.write("    - KeyName: %s\n" % key_name)
        outfile.write("      Value: %s\n" % prop)
        outfile.write("      CapcaseName: %s%s\n" % (setting_capcase_name, prop.title().replace('-', '')))
        outfile.write("      Type: %s\n" % value_type)
        constants[key_name]=prop
        if default_value is not None:
            outfile.write("      DefaultValue: %s\n" % default_value)

# generate enumerations
outfile.write("\nNMEnums:\n")
for enum in girxml.findall('./gi:namespace/gi:enumeration', ns_map):
    outfile.write("  - EnumClass: %s\n" % enum.attrib['name'])
    outfile.write("    Members:\n")
    for enumval in enum.findall('./gi:member', ns_map):
        cname = enumval.attrib['{%s}identifier' % ns_map['c']]
        cvalue = '%d' % int(enumval.attrib['value'])
        constants[cname]=cvalue
        outfile.write("    - Name: %s\n      Value: %s\n" % (cname, cvalue))

for enum in girxml.findall('./gi:namespace/gi:bitfield', ns_map):
    outfile.write("  - EnumClass: %s\n" % enum.attrib['name'])
    outfile.write("    Members:\n")
    for enumval in enum.findall('./gi:member', ns_map):
        cname = enumval.attrib['{%s}identifier' % ns_map['c']]
        cvalue = '0x%x' % int(enumval.attrib['value'])
        constants[cname]=cvalue
        outfile.write("    - Name: %s\n      Value: %s\n" % (cname, cvalue))

# generate gchar* constants (maybe are enumeration members, too)
outfile.write("  - EnumClass: StringConstants\n")
outfile.write("    Members:\n")
for const in girxml.findall('./gi:namespace/gi:constant', ns_map):
    cname = const.attrib['{%s}type' % ns_map['c']]
    cvalue = const.attrib['value']
    if cname not in constants and const.find('./gi:type[@c:type="gchar*"]', ns_map) is not None:
        constants[cname]=cvalue
        outfile.write("    - Name: %s\n      Value: \"%s\"\n" % (cname, cvalue))

outfile.close()
