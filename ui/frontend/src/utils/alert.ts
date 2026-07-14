import Swal from 'sweetalert2';

// 💡 성공 알림을 세련되게 띄워주는 공통 함수
export const ShowAPISuccess = (response: any, defaultMessage: string = "성공적으로 처리되었습니다.") => {
    // response가 Wails가 넘겨준 APIResponse 객체인지, 아니면 단순 문자열(message)인지 확인
    let successMessage = defaultMessage;

    if (response) {
        if (typeof response === 'string') {
            // 회원가입 등에서 성공 메시지 문자열만 넘어온 경우
            successMessage = response;
        } else if (response.message) {
            // 로그인 등에서 APIResponse[AuthResult] 구조체 객체가 넘어온 경우
            successMessage = response.message;
        }
    }

    // SweetAlert2 성공 팝업 출력
    Swal.fire({
        icon: 'success',
        title: successMessage,
        confirmButtonColor: '#3085d6'
    });
};


// 💡 모든 catch 에러를 넘겨받아 자동으로 언마샬링하고 SweetAlert를 띄워주는 함수
export const ShowAPIError = (error: any, defaultTitle: string = "오류 발생") => {
    // 에러 객체에서 순수 JSON 텍스트만 추출
    let errorStr = error?.message || String(error);
    if (errorStr.startsWith("Error: ")) {
        errorStr = errorStr.replace("Error: ", "");
    }

    try {
        // 언마샬링 시도
        const apiError = JSON.parse(errorStr);

        // 알맞은 자리에 데이터 바인딩
        Swal.fire({
            icon: 'error',
            title: apiError.message || defaultTitle,
            text: typeof apiError.details === 'object' ? JSON.stringify(apiError.details) : apiError.details,
            footer: `<span style="color:#aaa; font-size: 11px;">Error Location: ${apiError.code}</span>`
        });
    } catch {
        // JSON 구조가 아닌 시스템/네트워크 에러 대응
        Swal.fire({
            icon: 'error',
            title: defaultTitle,
            text: errorStr,
        });
    }
};
