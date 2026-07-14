import { useState, type ChangeEvent, type FormEvent } from "react";
import { ArrowRight, Building2, MessageSquareMore } from "lucide-react";
import { Link, useNavigate } from "react-router-dom";

import { ShowAPIError, ShowAPISuccess } from "../../utils/alert";

import { Register } from "../../../wailsjs/go/services/HTTPClientService"


type FormState = {
  login_id: string;
  name: string;
  password: string;
  confirmPassword: string;
};

const RegisterPage = () => {
  const navigate = useNavigate()

  const [formData, setFormData] = useState<FormState>({
    login_id: "",
    name: "",
    password: "",
    confirmPassword: "",
  });
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [isLoading, setIsLoading] = useState(false);

  const validateForm = () => {
    const nextErrors: Record<string, string> = {};

    if (!formData.login_id.trim()) {
      nextErrors.login_id = "로그인 아이디를 입력해주세요";
    }

    if (!formData.name.trim()) {
      nextErrors.name = "이름을 입력해주세요";
    }

    if (!formData.password.trim()) {
      nextErrors.password = "비밀번호를 입력해주세요";
    }

    if (!formData.confirmPassword.trim()) {
      nextErrors.confirmPassword = "비밀번호 재입력을 입력해주세요";
    } else if (formData.password !== formData.confirmPassword) {
      nextErrors.confirmPassword = "비밀번호가 일치하지 않습니다";
    }

    return nextErrors;
  };

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));

    if (errors[name]) {
      setErrors((prev) => {
        const nextErrors = { ...prev };
        delete nextErrors[name];
        return nextErrors;
      });
    }
  };

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const nextErrors = validateForm();
    setErrors(nextErrors);

    if (Object.keys(nextErrors).length > 0) {
      setIsLoading(false);
      return;
    }

    setIsLoading(true);
    try{
      const res = await Register({login_id: formData.login_id, password: formData.password, name: formData.name })
      ShowAPISuccess(res)

      navigate("/")
    } catch(error){
      ShowAPIError(error, "회원가입 실패")
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="w-full max-w-6xl px-4 py-8 sm:px-6 lg:px-8">
      <div className="overflow-hidden rounded-[32px] border border-slate-200/80 bg-white/90 shadow-[0_24px_80px_-24px_rgba(15,23,42,0.45)] backdrop-blur">
        <div className="grid lg:grid-cols-[0.95fr_1.05fr]">
          <div className="bg-gradient-to-br from-slate-700 via-slate-600 to-slate-500 p-8 text-white sm:p-10 lg:p-12">
            <div className="flex items-center gap-3">
              <div className="flex h-11 w-11 items-center justify-center rounded-2xl bg-white/15">
                <MessageSquareMore className="h-6 w-6" />
              </div>
              <div>
                <p className="text-sm text-slate-200">새 팀원 등록</p>
                <h2 className="text-lg font-semibold">사내 채팅 시작</h2>
              </div>
            </div>

            <div className="mt-10 space-y-4">
              <p className="text-sm font-medium uppercase tracking-[0.3em] text-slate-200">Join Team</p>
              <h3 className="text-3xl font-semibold leading-tight">업무 대화의 흐름을 바로 이어가세요</h3>
              <p className="max-w-md text-sm leading-7 text-slate-200">
                팀원 초대와 채널 구성까지 자연스럽게 이어질 수 있는 계정을 만들어보세요.
              </p>
            </div>

            <div className="mt-8 rounded-2xl border border-white/15 bg-white/10 p-4">
              <div className="flex items-center gap-2 text-sm text-slate-200">
                <Building2 className="h-4 w-4" />
                조직별 권한과 채널을 손쉽게 관리할 수 있어요
              </div>
            </div>
          </div>

          <div className="p-8 sm:p-10 lg:p-12">
            <div className="mb-8">
              <p className="text-sm font-semibold text-sky-600">계정 생성</p>
              <h1 className="mt-2 text-3xl font-semibold text-slate-900">팀 계정을 시작하세요</h1>
              <p className="mt-2 text-sm leading-6 text-slate-500">
                로그인 아이디, 이름, 비밀번호를 입력하고 비밀번호를 한 번 더 확인해주세요.
              </p>
            </div>

            <form className="space-y-4" onSubmit={handleSubmit}>
              <div className="space-y-2">
                <label htmlFor="login_id" className="block text-sm font-medium text-slate-700">
                  로그인 아이디
                </label>
                <input
                  id="login_id"
                  name="login_id"
                  type="text"
                  value={formData.login_id}
                  onChange={handleChange}
                  placeholder="jhkim"
                  className="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-900 outline-none transition focus:border-sky-500 focus:bg-white focus:ring-4 focus:ring-sky-100"
                />
                {errors.login_id && <p className="text-xs text-rose-500">{errors.login_id}</p>}
              </div>

              <div className="space-y-2">
                <label htmlFor="name" className="block text-sm font-medium text-slate-700">
                  이름
                </label>
                <input
                  id="name"
                  name="name"
                  type="text"
                  value={formData.name}
                  onChange={handleChange}
                  placeholder="홍길동"
                  className="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-900 outline-none transition focus:border-sky-500 focus:bg-white focus:ring-4 focus:ring-sky-100"
                />
                {errors.name && <p className="text-xs text-rose-500">{errors.name}</p>}
              </div>

              <div className="space-y-2">
                <label htmlFor="registerPassword" className="block text-sm font-medium text-slate-700">
                  비밀번호
                </label>
                <input
                  id="registerPassword"
                  name="password"
                  type="password"
                  value={formData.password}
                  onChange={handleChange}
                  placeholder="비밀번호를 입력하세요"
                  className="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-900 outline-none transition focus:border-sky-500 focus:bg-white focus:ring-4 focus:ring-sky-100"
                />
                {errors.password && <p className="text-xs text-rose-500">{errors.password}</p>}
              </div>

              <div className="space-y-2">
                <label htmlFor="confirmPassword" className="block text-sm font-medium text-slate-700">
                  비밀번호 재입력
                </label>
                <input
                  id="confirmPassword"
                  name="confirmPassword"
                  type="password"
                  value={formData.confirmPassword}
                  onChange={handleChange}
                  placeholder="비밀번호를 다시 입력하세요"
                  className="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-900 outline-none transition focus:border-sky-500 focus:bg-white focus:ring-4 focus:ring-sky-100"
                />
                {errors.confirmPassword && <p className="text-xs text-rose-500">{errors.confirmPassword}</p>}
              </div>

              <button
                type="submit"
                disabled={isLoading}
                className="inline-flex w-full items-center justify-center gap-2 rounded-2xl bg-slate-900 px-4 py-3 text-sm font-semibold text-white transition hover:bg-slate-700 disabled:cursor-not-allowed disabled:bg-slate-400"
              >
                {isLoading ? "가입 중..." : "계정 생성하기"}
                <ArrowRight className="h-4 w-4" />
              </button>
            </form>

            <div className="mt-6 text-center text-sm text-slate-500">
              이미 계정이 있으신가요?{" "}
              <Link to="/" className="font-semibold text-sky-600 transition hover:text-sky-700">
                로그인하기
              </Link>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default RegisterPage;