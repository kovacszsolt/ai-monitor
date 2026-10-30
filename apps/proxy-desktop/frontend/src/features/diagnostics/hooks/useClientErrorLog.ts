import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { clearClientErrors, listClientErrors } from "@/lib/clientErrorCapture";

export function useClientErrorLog(page: number, pageSize: number) {
  const queryClient = useQueryClient();

  const listQ = useQuery({
    queryKey: ["client-errors", page, pageSize],
    queryFn: () => listClientErrors(page, pageSize),
  });

  const clearMut = useMutation({
    mutationFn: clearClientErrors,
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["client-errors"] });
    },
  });

  const totalPages = listQ.data ? Math.max(1, Math.ceil(listQ.data.total / pageSize)) : 1;

  return { listQ, clearMut, totalPages };
}
