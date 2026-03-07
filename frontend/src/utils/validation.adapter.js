import { z } from 'zod'

export function createFormRules(schema, path = '') {

    const rules = {}

    // Handle Objects
    if (schema instanceof z.ZodObject) {
        const shape = schema.shape
        for (const key in shape) {
            const fieldPath = path ? `${path}.${key}` : key
            Object.assign(rules, createFormRules(shape[key], fieldPath))
        }
    }

    // Create the validator for the current path
    if (path) {
        rules[path] = [{
            validator: (rule, value) => {
                const result = schema.safeParse(value)
                if (!result.success) {
                    const errMessage = result.error?.issues.map(err => {
                        return err.message
                    }).join('\n')
                    return new Error(errMessage)
                }
                return true
            },
            trigger: ['input', 'blur']
        }]
    }

    return rules
}