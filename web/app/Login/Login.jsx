import React from 'react';
import { connect } from 'react-redux';
import { login } from './duck';

const Login = ({ onSubmit }) => {
  let inputLogin;
  let inputPassword;

  const submit = (e) => {
    e.preventDefault();

    const username = inputLogin.value.trim();
    const password = inputPassword.value.trim();

    if (!login || !password) {
      return;
    }

    onSubmit({
      login: username,
      password
    });
  };

  return (
    <div className="container">
      <div className="row justify-content-md-center">
        <div className="col col-lg-4">
          <h3
            style={{
              textAlign: 'center',
              margin: '20px auto'
            }}
          >
          Sign in to Pulpe
        </h3>
          <div className="card">
            <div className="card-block">
              <form onSubmit={submit}>
                <div className="form-group">
                  <label htmlFor="loginField">Login or email address</label>
                  <input
                    type="text"
                    id="loginField"
                    className="form-control"
                    required
                    autoFocus
                    ref={(node) => { inputLogin = node; }}
                  />
                </div>
                <div className="form-group">
                  <label htmlFor="passwordField">Password</label>
                  <input
                    type="password"
                    id="passwordField"
                    className="form-control"
                    required
                    ref={(node) => { inputPassword = node; }}
                  />
                </div>
                <button className="btn btn-primary btn-block" type="submit">Sign in</button>
              </form>
            </div>
          </div>
          <div className="card" style={{ textAlign: 'center' }}>
            <div className="card-block">
              <p className="card-text">New to Pulpe? <a href="/join">Create an account</a>.</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default connect(
  null,
  {
    onSubmit: login
  }
)(Login);
