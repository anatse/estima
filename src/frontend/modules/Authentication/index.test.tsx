import * as React from 'react';
import {shallow} from 'enzyme';

import {Authentication} from './index';

describe('<Authentication /> тестирование страницы авторизации.', () => {
    const authenticationShallow = shallow(<Authentication />);

    const userLogin = 'user123';
    const userPassword = 'password123';

    const loginIdElement = '#Authentication__input__login';
    const passwordIdElement = '#Authentication__input__password';
    const enterIdElement = '#Authentication__button__enter';

    test('Есть два поля для ввода.', () => {
        expect(authenticationShallow.find(passwordIdElement)).toHaveLength(1);
        expect(authenticationShallow.find(loginIdElement)).toHaveLength(1);
    });

    test('Есть кнопка для авторизации.', () => {
        expect(authenticationShallow.find(enterIdElement)).toHaveLength(1);
    });

    test(`Ввод логина: ${userLogin}`, () => {
        authenticationShallow
            .find(loginIdElement)
            .first()
            .simulate('change', { target: { value: userLogin } });
        expect(authenticationShallow
            .find(loginIdElement)
            .first()
            .props()
            .value)
            .toEqual(userLogin);
    });

    test(`Ввод пароля: ${userPassword}`, () => {
        authenticationShallow
            .find(passwordIdElement)
            .last()
            .simulate('change', { target: { value: userPassword } });
        expect(authenticationShallow
            .find(passwordIdElement)
            .last()
            .props().value)
            .toEqual(userPassword);
    });

    test(`Метод handlerSignIn вызвал запрос с двумя параметрами (${userLogin}, ${userPassword}).`, () => {
        const inst = authenticationShallow.instance() as any;
        inst.getAuthorization = jest.fn();
        inst.handlerSignIn();

        expect(inst.getAuthorization).toHaveBeenCalledTimes(1);
        expect(inst.getAuthorization).toHaveBeenCalledWith(userLogin, userPassword);
    });

});
