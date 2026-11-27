// Package system 提供面板所在机器的系统能力（如文件选择对话框）。
package system

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// fileDialogScript 是 Windows 文件选择对话框脚本模板。
// 不用 WinForms OpenFileDialog（样式老旧），而是嵌入 C# 直接调用 COM
// IFileOpenDialog —— 资源管理器原生现代对话框。
// 首行按项目约定统一 UTF-8，保证中文路径输出不乱码。
// %s 依次为对话框标题、文件过滤器，均为 PowerShell 单引号字符串字面量（已转义）。
const fileDialogScript = `$OutputEncoding = [console]::InputEncoding = [console]::OutputEncoding = [Text.UTF8Encoding]::UTF8
$cs = @'
using System;
using System.Runtime.InteropServices;

public static class ModernFileDialog
{
    [ComImport, Guid("42F85136-DB7E-439C-85F1-E4075D135FC8"), InterfaceType(ComInterfaceType.InterfaceIsIUnknown)]
    private interface IFileOpenDialog
    {
        [PreserveSig] int Show(IntPtr hwndOwner);
        void SetFileTypes(uint cFileTypes, [MarshalAs(UnmanagedType.LPArray)] COMDLG_FILTERSPEC[] rgFilterSpec);
        void SetFileTypeIndex(uint iFileType);
        void GetFileTypeIndex(out uint piFileType);
        void Advise(IntPtr pfde, out uint pdwCookie);
        void Unadvise(uint dwCookie);
        void SetOptions(uint fos);
        void GetOptions(out uint pfos);
        void SetDefaultFolder(IShellItem psi);
        void SetFolder(IShellItem psi);
        void GetFolder(out IShellItem ppsi);
        void GetCurrentSelection(out IShellItem ppsi);
        void SetFileName([MarshalAs(UnmanagedType.LPWStr)] string pszName);
        void GetFileName([MarshalAs(UnmanagedType.LPWStr)] out string pszName);
        void SetTitle([MarshalAs(UnmanagedType.LPWStr)] string pszTitle);
        void SetOkButtonLabel([MarshalAs(UnmanagedType.LPWStr)] string pszText);
        void SetFileNameLabel([MarshalAs(UnmanagedType.LPWStr)] string pszLabel);
        void GetResult(out IShellItem ppsi);
        void AddPlace(IShellItem psi, int fdap);
        void SetDefaultExtension([MarshalAs(UnmanagedType.LPWStr)] string pszDefaultExtension);
        void Close(int hr);
        void SetClientGuid(ref Guid guid);
        void ClearClientData();
        void SetFilter(IntPtr pFilter);
        void GetResults(out IntPtr ppenum);
        void GetSelectedItems(out IntPtr ppsai);
    }

    [ComImport, Guid("43826D1E-E718-42EE-BC55-A1E261C37BFE"), InterfaceType(ComInterfaceType.InterfaceIsIUnknown)]
    private interface IShellItem
    {
        void BindToHandler(IntPtr pbc, ref Guid bhid, ref Guid riid, out IntPtr ppv);
        void GetParent(out IShellItem ppsi);
        void GetDisplayName(uint sigdnName, out IntPtr ppszName);
        void GetAttributes(uint sfgaoMask, out uint psfgaoAttribs);
        void Compare(IShellItem psi, uint hint, out int piOrder);
    }

    [StructLayout(LayoutKind.Sequential, CharSet = CharSet.Unicode)]
    private struct COMDLG_FILTERSPEC
    {
        [MarshalAs(UnmanagedType.LPWStr)] public string pszName;
        [MarshalAs(UnmanagedType.LPWStr)] public string pszSpec;
    }

    private const int S_OK = 0;
    private const uint FOS_FORCEFILESYSTEM = 0x40;
    private const uint FOS_FILEMUSTEXIST = 0x1000;
    private const uint SIGDN_FILESYSPATH = 0x80058000;

    // filter 为 WinForms 风格 "名称 (*.jar)|*.jar|名称 (*.*)|*.*"
    public static string PickFile(string title, string filter)
    {
        var dlg = (IFileOpenDialog)Activator.CreateInstance(
            Type.GetTypeFromCLSID(new Guid("DC1C5A9C-E88A-4DDE-A5A1-60F82A20AEF7")));
        try
        {
            dlg.SetOptions(FOS_FORCEFILESYSTEM | FOS_FILEMUSTEXIST);
            if (!string.IsNullOrEmpty(title)) dlg.SetTitle(title);
            var specs = ParseFilter(filter);
            if (specs.Length > 0) dlg.SetFileTypes((uint)specs.Length, specs);
            if (dlg.Show(IntPtr.Zero) != S_OK) return null; // 用户取消
            // 注意：PS 5.1 的 Add-Type 默认 C# 5 编译器，不支持内联 out 变量
            IShellItem item;
            dlg.GetResult(out item);
            IntPtr ppsz;
            item.GetDisplayName(SIGDN_FILESYSPATH, out ppsz);
            try { return Marshal.PtrToStringUni(ppsz); }
            finally { Marshal.FreeCoTaskMem(ppsz); }
        }
        finally { Marshal.FinalReleaseComObject(dlg); }
    }

    private static COMDLG_FILTERSPEC[] ParseFilter(string filter)
    {
        if (string.IsNullOrEmpty(filter)) return new COMDLG_FILTERSPEC[0];
        var parts = filter.Split('|');
        var list = new System.Collections.Generic.List<COMDLG_FILTERSPEC>();
        for (int i = 0; i + 1 < parts.Length; i += 2)
            list.Add(new COMDLG_FILTERSPEC { pszName = parts[i], pszSpec = parts[i + 1] });
        return list.ToArray();
    }
}
'@
Add-Type -TypeDefinition $cs
$path = [ModernFileDialog]::PickFile('%s', '%s')
if ($path) { [Console]::Out.Write($path) }
`

// PickFile 在面板所在机器上弹出系统文件选择对话框，返回选中文件的完整路径。
// 用户取消对话框时返回空字符串（不视为错误）。
// filter 形如 "服务端核心 (*.jar)|*.jar|所有文件 (*.*)|*.*"；为空时使用默认 jar 过滤器。
// 仅 Windows 可用；其他平台返回错误。
func PickFile(ctx context.Context, title, filter string) (string, error) {
	if runtime.GOOS != "windows" {
		return "", fmt.Errorf("当前平台 %s 暂不支持文件选择对话框", runtime.GOOS)
	}
	if title == "" {
		title = "选择文件"
	}
	if filter == "" {
		filter = "服务端核心 (*.jar)|*.jar|所有文件 (*.*)|*.*"
	}

	script := fmt.Sprintf(fileDialogScript, escapePSString(title), escapePSString(filter))
	// -EncodedCommand 要求 Base64(UTF-16LE)，避免脚本中的中文经命令行传参时乱码
	cmd := exec.CommandContext(ctx, "powershell", "-NoProfile", "-NonInteractive", "-EncodedCommand", encodedCommand(script))
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("打开文件选择对话框: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// escapePSString 转义 PowerShell 单引号字符串字面量（单引号双写）。
func escapePSString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
