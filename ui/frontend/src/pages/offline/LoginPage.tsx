import { useState, type ChangeEvent, type FormEvent } from "react";
import { ArrowRight, LockKeyhole, MessageSquareMore, ReceiptSwissFranc, ShieldCheck } from "lucide-react";
import { Link, useNavigate } from "react-router-dom";

import { ShowAPIError, ShowAPISuccess } from "../../utils/alert";
import useAuthStore from "../../stores/authStore";

import { ConnectWebSocket, Login } from "../../../wailsjs/go/services/HTTPClientService"

type FormState = {
  login_id: string;
  password: string;
};

const LoginPage = () => {
  const navigate = useNavigate()
  const setAuth = useAuthStore((state) => state.setAuth)

  const [formData, setFormData] = useState<FormState>({
    login_id: "",
    password: "",
  });
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [isLoading, setIsLoading] = useState(false);

  const validateForm = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.login_id.trim()) {
      newErrors.login_id = "아이디를 입력해주세요";
    }

    if (!formData.password.trim()) {
      newErrors.password = "비밀번호를 입력해주세요";
    }

    return newErrors;
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
    const newErrors = validateForm();
    setErrors(newErrors);

    if (Object.keys(newErrors).length > 0) {
      setIsLoading(false);
      return;
    }

    setIsLoading(true);
    try {
      const res = await Login({login_id: formData.login_id, password: formData.password })
      ShowAPISuccess(res)

      const user = res?.data?.user
      const token = res?.data?.token

      if (user && token) {
        setAuth(user, token)
        await ConnectWebSocket()
      }

      navigate("/online/friends")
    } catch (error) {
      ShowAPIError(error, "로그인 실패")
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="w-full max-w-6xl px-4 py-8 sm:px-6 lg:px-8">
      <div className="overflow-hidden rounded-[32px] border border-slate-200/80 bg-white/90 shadow-[0_24px_80px_-24px_rgba(15,23,42,0.45)] backdrop-blur">
        <div className="grid lg:grid-cols-[1.02fr_0.98fr]">
          <div className="bg-gradient-to-br from-slate-900 via-slate-800 to-slate-700 p-8 text-white sm:p-10 lg:p-12">
            <div className="flex items-center gap-3">
              <div className="flex h-11 w-11 items-center justify-center rounded-2xl bg-white/15">
                <MessageSquareMore className="h-6 w-6" />
              </div>
              <div>
                <p className="text-sm text-slate-300">사내 업무용</p>
                <h2 className="text-lg font-semibold">G-kk-ch</h2>
              </div>
            </div>

            <div className="mt-10 space-y-4">
              <p className="text-sm font-medium uppercase tracking-[0.3em] text-sky-300">Internal ZomPPang Chat</p>
              <h3 className="text-3xl font-semibold leading-tight">좀팽이 튀기기 21버전</h3>
              <p className="max-w-md text-sm leading-7 text-slate-300">
                {/* 공지, 회의, 업무 공유를 한 번에 이어가고 실시간으로 확인해보세요. */}
              </p>
            </div>

            <div className="mt-8 rounded-2xl border border-white/15 bg-white/10 p-4">
              <div className="flex items-center gap-2 text-sm text-slate-200">
                <ShieldCheck className="h-4 w-4 text-sky-300" />
                회사 도메인 기반 인증으로 안전하게
              </div>
              <div className="mt-3 flex items-center gap-2 text-sm text-slate-300">
                <LockKeyhole className="h-4 w-4 text-sky-300" />
                민감한 업무 내용도 안심하고 관리하세요
              </div>
            </div>
          </div>

          <div className="p-8 sm:p-10 lg:p-12">
            <div className="mb-8">
              <p className="text-sm font-semibold text-sky-600">로그인</p>
              <h1 className="mt-2 text-3xl font-semibold text-slate-900">환영합니다</h1>
              <p className="mt-2 text-sm leading-6 text-slate-500">
                사내 채팅 계정으로 로그인해 주세요.
              </p>
            </div>

            <form className="space-y-5" onSubmit={handleSubmit}>
              <div className="space-y-2">
                <label htmlFor="login_id" className="block text-sm font-medium text-slate-700">
                  아이디
                </label>
                <input
                  type="text"
                  id="login_id"
                  name="login_id"
                  value={formData.login_id}
                  onChange={handleChange}
                  placeholder="예: jhkim"
                  className="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-900 outline-none transition focus:border-sky-500 focus:bg-white focus:ring-4 focus:ring-sky-100"
                />
                {errors.login_id && <p className="text-xs text-rose-500">{errors.login_id}</p>}
              </div>

              <div className="space-y-2">
                <label htmlFor="password" className="block text-sm font-medium text-slate-700">
                  비밀번호
                </label>
                <input
                  type="password"
                  id="password"
                  name="password"
                  value={formData.password}
                  onChange={handleChange}
                  placeholder="비밀번호를 입력하세요"
                  className="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-900 outline-none transition focus:border-sky-500 focus:bg-white focus:ring-4 focus:ring-sky-100"
                />
                {errors.password && <p className="text-xs text-rose-500">{errors.password}</p>}
              </div>

              <button
                type="submit"
                disabled={isLoading}
                className="inline-flex w-full items-center justify-center gap-2 rounded-2xl bg-slate-900 px-4 py-3 text-sm font-semibold text-white transition hover:bg-slate-700 disabled:cursor-not-allowed disabled:bg-slate-400"
              >
                {isLoading ? "로그인 중..." : "로그인"}
                <ArrowRight className="h-4 w-4" />
              </button>
            </form>

            <div className="mt-6 text-center text-sm text-slate-500">
              계정이 아직 없으신가요?{" "}
              <Link to="/register" className="font-semibold text-sky-600 transition hover:text-sky-700">
                계정 생성하기
              </Link>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default LoginPage;