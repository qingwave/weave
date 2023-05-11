import Mock, { mock } from 'mockjs';
import data from './data';

for (let item of data) {
    Mock.mock(item.path, item.method, item.do ? item.do : (params) => {
        console.log("Mock ->", params)
        return {
            code: 200,
            data: item.data,
        }
    })
}

export default Mock;
