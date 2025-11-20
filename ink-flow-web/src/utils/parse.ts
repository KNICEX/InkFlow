import type { LocationQueryValue, RouteRecordNameGeneric } from 'vue-router'

export const parseRouteParam = (param: string | string[] | undefined) => {
  if (param == undefined) {
    return ''
  }
  if (typeof param === 'string') {
    return param
  } else if (param.length === 0) {
    return ''
  }
  return param[0]
}

export const parseRouteParamToInt = (param: string | string[] | undefined) => {
  if (param == undefined) {
    return 0
  }
  if (typeof param === 'string') {
    return parseInt(param == '' ? '0' : param)
  } else if (param.length === 0) {
    return 0
  }
  return parseInt(param[0] == '' ? '0' : param[0])
}

export const parseRouteQuery = (query: LocationQueryValue | LocationQueryValue[]): string => {
  if (query == undefined) {
    return ''
  }
  if (typeof query === 'string') {
    return query
  }
  return query[0] ?? ''
}

export const parseResponseDate = (data: Record<string, any>) => {
  if (!data) {
    return
  }
  if (Array.isArray(data)) {
    for (const item of data) {
      parseResponseDate(item)
    }
  } else if (typeof data === 'object') {
    // 检查属性
    for (const key in data) {
      if (typeof data[key] === 'object') {
        parseResponseDate(data[key])
      }
    }

    // 处理日期字段
    const dateFields = ['createdAt', 'updatedAt', 'date', 'time', 'timestamp']
    for (const field of dateFields) {
      if (data[field]) {
        try {
          // 尝试直接解析
          const date = new Date(data[field])
          // 检查是否为有效日期
          if (!isNaN(date.getTime())) {
            data[field] = date
          } else {
            // 尝试解析ISO格式
            const isoDate = new Date(data[field].replace(' ', 'T'))
            if (!isNaN(isoDate.getTime())) {
              data[field] = isoDate
            }
          }
        } catch (e) {
          console.warn(`Failed to parse date for field ${field}:`, e)
        }
      }
    }
  }
}
