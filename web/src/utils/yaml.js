
import YAML from 'yaml'
import yaml from 'js-yaml'

export const json2yaml = (json) => {
    try {
        return yaml.dump(JSON.parse(json))
    } catch (err) {
        console.log("err convert json to yaml: ", err)
        return ''
    }
}

export const yaml2json = (yamlstr) => {
    try {
        return JSON.stringify(YAML.parse(yamlstr), null, 2)
    } catch (err) {
        console.log("err convert yaml to json: ", err)
        return ''
    }
}

export const obj2yaml = (obj) => {
    try {
        return yaml.dump(obj)
    } catch (err) {
        console.log("err convert obj to yaml: ", err)
        return ''
    }
}

export const yaml2obj = (yamlstr) => {
    try {
        return YAML.parse(yamlstr)
    } catch (err) {
        console.log("err convert yaml to obj: ", err)
        return ''
    }
}
