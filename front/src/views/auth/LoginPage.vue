<script setup>
import "@/styles/main.css";
</script>

<template>
  <div class="login">
    <div class="login__container">
      <div class="login_in">
        <div class="login_header">
          <a href="login1.html">
            <img src="@/assets/img/svgs/arrow-left.svg" alt="" />
          </a>
        </div>
        <div class="login_body">
          <div class="login_body_in">
            <h6 class="text_32 text_white text_bold pb-5">Регистрация</h6>
            <div class="form-group my_form-group">
              <div class="my_input">
                <input type="text" class="form-control" placeholder="Имя" />
              </div>
              <div class="my_input">
                <input type="text" class="form-control" placeholder="Фамилия" />
              </div>
              <div class="my_input">
                <input type="tel" class="form-control" placeholder="Телефон" />
              </div>
              <div class="my_input">
                <input
                  type="password"
                  class="form-control"
                  placeholder="Придумайте пароль"
                />
              </div>
            </div>

            <div class="my_check">
              <input class="form-check-input" type="checkbox" id="flexCheck" />
              <label class="form-check-label .text_semi_bold" for="flexCheck">
                <span>Принимаю</span> условия пользования
              </label>
            </div>
 
            <button type="button" class="btn blue_btn">
              Зарегистрироваться
            </button>
          </div>
        </div>
        <div class="login_footer">
          <div class="login_footer__left">
            <p class="login_footer_p">Выберите язык:</p>
            <div class="login_footer__left__lang">
              <button class="active">Русский</button>
              <button>Ўзбек</button>
            </div>
          </div>
          <div class="login_footer__right">
            <p class="login_footer_p">© 2022 Genesis Innovation</p>
            <button class="btn login_footer__right_blur_btn">
              Связаться с нами
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style>
.login {
  background-color: #011564;
  background-image: url("@/assets/img/images/");
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  height: 100vh;
  overflow: auto;
}

.login__container {
  margin: 0 auto;
  padding: 0 40px;
  height: 100%;
}

.login_header {
  padding: 40px 0;
}

.login_in {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  height: 100%;
}

.login_body_p {
  font-size: 14px;
  line-height: 1.2;
  color: rgba(255, 255, 255, 0.6);
  font-family: "Inter-semiBold", sans-serif;
  min-height: 50px;
}

.login_body {
  flex: auto;
  max-width: 900px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.login_body_in {
  max-width: 310px;
  margin: 0 auto;
}

.login .my_input input {
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  margin-bottom: 12px;
  min-height: 52px;
  font-family: "Inter-semiBold", sans-serif;
  font-size: 14px;
  color: #ffffff;
}

.login .my_input input:focus {
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.6);
  outline: none;
  box-shadow: none;
}

.login .my_input.error > input::placeholder {
  color: #f23d5b;
  opacity: 1;
}

.login .my_input.error > input {
  background: rgba(242, 61, 91, 0.12);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

/*checked*/
.my_check {
  padding: 25px 50px 32px 0;
  position: relative;
  left: 20px;
}

.my_check input {
  position: absolute;
  width: 16px;
  height: 16px;
}

.my_check label {
  color: rgb(255, 255, 255, 0.8);
  font-size: 12px;
  padding: 0 50px 0 0;
  font-family: "Inter-semiBold", sans-serif;
}

.my_check label span {
  color: rgb(255, 255, 255, 0.4);
}

.login_body_link {
  font-family: "Inter-semiBold", sans-serif;
  font-size: 12px;
  line-height: 1.2;
  color: rgba(255, 255, 255, 0.8);
  margin-bottom: 32px;
  display: inline-block;
}

.login_body_link:hover,
.login_body_link2:hover {
  color: rgba(255, 255, 255, 0.8);
}

.login .blue_btn {
  width: 100%;
  height: 52px;
  background-color: #1d65f5;
  border-radius: 8px;
  font-family: "Inter-bold", sans-serif;
  font-size: 16px;
  line-height: 1.2;
  color: #ffffff;
}

.login_body_link2 {
  font-family: "Inter-semiBold", sans-serif;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.6);
}

.login_footer {
  padding: 40px 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.login_footer__left {
  display: flex;
  align-items: center;
}

.login_footer_p {
  font-family: "Inter-semiBold", sans-serif;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.6);
  margin-right: 10px;
  margin-bottom: 0;
}

.login_footer__left__lang button {
  font-family: "Inter-semiBold", sans-serif;
  font-size: 14px;
  color: #ffffff;
  padding: 4px 8px;
  background-color: transparent;
  border: 1px solid transparent;
  border-radius: 8px;
}

.login_footer__left__lang button.active {
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.login_footer__right {
  display: flex;
  align-items: center;
}

.login_footer__right_blur_btn {
  padding: 4px 8px;
  background-color: #1d65f5;
  border-radius: 8px;
  font-weight: 600;
  font-size: 14px;
  color: #ffffff;
  font-family: "Inter-semiBold", sans-serif;
}

.login_footer__right_blur_btn:hover {
  color: #fff;
}

@media only screen and (max-width: 768px) {
  .login_footer {
    flex-direction: column;
  }

  .login_footer__left {
    margin-bottom: 12px;
  }

  .login__container {
    padding: 0 16px;
  }

  .login_header,
  .login_footer {
    padding: 20px 0;
  }

  .my_check {
    padding: 0;
  }
}
</style>
