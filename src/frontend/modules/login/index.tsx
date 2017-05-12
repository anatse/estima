import * as React from 'react';
import './styles/login.css';
import {Redirect} from "react-router";

interface IState {
    login: string;
    password: string;
    isLoading: boolean;
    isLogged: false;
}

/**
 * Страница авторизации пользователя.
 */
export class Login extends React.Component<void, IState> {

    state: IState = {
        login: '',
        password: '',
        isLoading: false,
        isLogged: false
    };

    componentWillMount () {
        document.title = 'Логин';
    }

    /**
     * Авторизация в системе.
     */
    handlerSignIn = () => {
        const {
            login,
            password,
        } = this.state;

        this.setState(() => {
            return {isLoading: true};
        });

        var url = '/login';
        fetch(url, {
            method: "post",
            body: JSON.stringify({uname: login, upass: password}),
            headers: {
                'Accept': 'application/json',
                'Cache': 'no-cache'
            },
            credentials: "same-origin"
        }).then((response) => {
            response.text().then ((text) => {
                this.setState(() => {
                    return {
                        isLoading: false,
                        isLogged: true
                    };
                });

                // console.log (text)
            })
        })
    }

    render() {
        const {
            login,
            password,
            isLoading,
            isLogged
        } = this.state;

        if (isLogged) {
            return <Redirect to="/projects"/>
        }

        return (
            <div className="login">
                <form className="login_form">
                    <div className="login_row login_row__align_center">
                        <h1 className="login_header">Estimator</h1>
                    </div>
                    <div className="login_row login_row__align_center">
                        <input
                            className="login_input"
                            type="text"
                            placeholder="Логин"
                            onChange={event => this.setState({login: event.target.value})}
                            value={login}
                        />
                    </div>
                    <div className="login_row login_row__align_center">
                        <input
                            className="login_input"
                            type="password"
                            placeholder="Пароль"
                            onChange={event => this.setState({password: event.target.value})}
                            value={password}
                        />
                    </div>
                    <div className="login_row login_row__align_center">
                        <button
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
