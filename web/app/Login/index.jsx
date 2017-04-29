import React from 'react';

const Login = () =>
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
            <form>
              <div className="form-group">
                <label htmlFor="loginField">Username or email address</label>
                <input type="text" id="loginField" className="form-control" required autoFocus />
              </div>
              <div className="form-group">
                <label htmlFor="passwordField">Password</label>
                <input type="password" id="passwordField" className="form-control" required />
              </div>
              <button className="btn btn-primary btn-block" type="submit">Sign in</button>
            </form>
          </div>
        </div>

      </div>

    </div>
  </div>;

export default Login;
