import React from 'react';
import { connect } from 'react-redux';
import { register } from './duck';

const Register = ({ onSubmit }) => {
  let inputEmail;
  let inputFullName;
  let inputPassword;

  const submit = (e) => {
    e.preventDefault();

    const email = inputEmail.value.trim();
    const fullName = inputFullName.value.trim();
    const password = inputPassword.value.trim();

    if (!email || !fullName || !password) {
      return;
    }

    onSubmit({
      email,
      fullName,
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
          Sign up to Pulpe
        </h3>
          <div className="card">
            <div className="card-block">
              <form onSubmit={submit}>
                <div className="form-group">
                  <label htmlFor="emailField">Email address</label>
                  <input
                    type="email"
                    id="emailField"
                    className="form-control"
                    required
                    autoFocus
                    ref={(node) => { inputEmail = node; }}
                  />
                </div>
                <div className="form-group">
                  <label htmlFor="fullNameField">Full name</label>
                  <input
                    type="text"
                    id="fullNameField"
                    className="form-control"
                    required
                    ref={(node) => { inputFullName = node; }}
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
                <button className="btn btn-primary btn-block" type="submit">Sign up</button>
              </form>
            </div>
          </div>
          <div className="card" style={{ textAlign: 'center' }}>
            <div className="card-block">
              <p className="card-text">Already registered? <a href="/login">Sign in here</a>.</p>
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
    onSubmit: register
  }
)(Register);
