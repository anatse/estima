import * as React from 'react';
import './styles/login.css';

interface IState {
    login: string;
    password: string;
    isLoading: boolean;
}

/**
 * Страница авторизации пользователя.
 */
export class Login extends React.Component<void, IState> {

    state: IState = {
        login: '',
        password: '',
        isLoading: false,
    }

    /**
     * Авторизация в системе.
     */
    handlerSignIn = () => {
        const {
            login,
            password,
        } = this.state;

        console.log('start login with: ', login, password);
    }

    render() {
        const {
            login,
            password,
            isLoading,
        } = this.state;

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
