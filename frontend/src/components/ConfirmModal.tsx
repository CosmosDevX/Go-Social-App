interface ConfirmModalProps {
  open: boolean
  title: string
  message: string
  confirmLabel?: string
  cancelLabel?: string
  loading?: boolean
  danger?: boolean
  onConfirm: () => void
  onCancel: () => void
}

export function ConfirmModal({
  open,
  title,
  message,
  confirmLabel = 'Удалить',
  cancelLabel = 'Отмена',
  loading = false,
  danger = true,
  onConfirm,
  onCancel,
}: ConfirmModalProps) {
  if (!open) return null

  return (
    <div
      className="fixed inset-0 z-[100] flex items-center justify-center p-4"
      role="dialog"
      aria-modal="true"
    >
      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-black/70 backdrop-blur-sm"
        onClick={loading ? undefined : onCancel}
      />

      {/* Panel */}
      <div className="relative w-full max-w-sm glass-strong p-6 shadow-2xl shadow-black/40 animate-[fadeIn_0.15s_ease-out]">
        <h3 className="text-lg font-semibold text-white mb-2">{title}</h3>
        <p className="text-sm text-white/60 mb-6 leading-relaxed">{message}</p>

        <div className="flex gap-3 justify-end">
          <button
            type="button"
            className="btn-ghost text-sm px-4 py-2"
            onClick={onCancel}
            disabled={loading}
          >
            {cancelLabel}
          </button>
          <button
            type="button"
            onClick={onConfirm}
            disabled={loading}
            className={`
              text-sm px-4 py-2 rounded-xl font-medium transition-all duration-200
              disabled:opacity-50 disabled:cursor-not-allowed active:scale-[0.98]
              ${danger
                ? 'bg-red-500/90 hover:bg-red-500 text-white shadow-lg shadow-red-500/20'
                : 'btn-primary'
              }
            `}
          >
            {loading ? 'Удаление...' : confirmLabel}
          </button>
        </div>
      </div>
    </div>
  )
}
