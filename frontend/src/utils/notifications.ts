type NotificationType = "success" | "error" | "info";

export const showNotification = (
  message: string,
  type: NotificationType = "info"
) => {
  // For now, using alert. In a real app, you'd use a proper notification library
  // like react-hot-toast, react-toastify, etc.
  switch (type) {
    case "success":
      alert(`✅ ${message}`);
      break;
    case "error":
      alert(`❌ ${message}`);
      break;
    case "info":
    default:
      alert(`ℹ️ ${message}`);
      break;
  }
};

export const notifications = {
  registerSuccess: () =>
    showNotification(
      "Registration successful, you will redirect to login page",
      "success"
    ),
  logoutSuccess: () => showNotification("Logged out successfully", "success"),
  taskCreated: () => showNotification("Task created successfully", "success"),
  taskUpdated: () => showNotification("Task updated successfully", "success"),
  taskDeleted: () => showNotification("Task deleted successfully", "success"),
  taskCreateFailed: (error: string) =>
    showNotification(`Failed to create task: ${error}`, "error"),
  taskUpdateFailed: (error: string) =>
    showNotification(`Failed to update task: ${error}`, "error"),
  taskDeleteFailed: (error: string) =>
    showNotification(`Failed to delete task: ${error}`, "error"),
};
