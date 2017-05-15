import * as React from 'react';
import {Redirect} from 'react-router';

import './styles/login.css';

interface IState {
    login: string;
    password: string;
    isLoading: boolean;
    isLogged: false;
}

/**
 * Страница авторизации пользователя.
 */
export class Authentication extends React.Component<any, IState> {

    state: IState = {
        login: '',
        password: '',
        isLoading: false,
        isLogged: false,
    };

    componentWillMount () {
        document.title = 'Авторизация';
    }

    /**
     * Вызов сервиса для авторизации пользователя в системе.
     *
     * @param login {string} Логин пользователя.
     * @param password {string} Пароль пользователя.
     */
    getAuthorization (login: string, password: string) {
        this.setState({isLoading: true});
        const url = '/login';
        fetch(url, {
            method: 'post',
            body: JSON.stringify({uname: login, upass: password}),
            headers: {
                Accept: 'application/json',
                Cache: 'no-cache',
            },
            credentials: 'same-origin',
        }).then((response) => {
            response.json().then ((response) => {
                this.setState(() => {
                    return {
                        isLoading: false,
                        isLogged: true,
                    };
                });

                console.log ('response', response.body);
            }, (error) => {
                console.log ('error json', error);
            });
        }, (error) => {
            console.log ('error connect', error);
        });
    }

    /**
     * Авторизация в системе.
     */
    handlerSignIn = () => {
        const {
            login,
            password,
        } = this.state;

        this.getAuthorization(login, password);
    }

    render() {
        const {
            login,
            password,
            isLoading,
            isLogged,
        } = this.state;

        if (isLogged) {
            return <Redirect to="/projects" />;
        }

        return (
            <div className="login">
                <form className="login_form">
                    <div className="login_row login_row__align_center">
                        <h1 className="login_header">Estimator</h1>
                    </div>
                    <div className="login_row login_row__align_center">
                        <input
                            id="Authentication__input__login"
                            className="login_input"
                            type="text"
                            placeholder="Логин"
                            onChange={event => this.setState({login: event.target.value})}
                            value={login}
                        />
                    </div>
                    <div className="login_row login_row__align_center">
                        <input
                            id="Authentication__input__password"
                            className="login_input"
                            type="password"
                            placeholder="Пароль"
                            onChange={event => this.setState({password: event.target.value})}
                            value={password}
                        />
                    </div>
                    <div className="login_row login_row__align_center">
                        <button
                            id="Authentication__button__enter"
                            type="button"
                            disabled={isLoading}
                            className="login_sign-in"
                            onClick={this.handlerSignIn}
                        >
                            Вход
                        </button>
                    </div>
                </form>
            </div>
        );
    }
}
