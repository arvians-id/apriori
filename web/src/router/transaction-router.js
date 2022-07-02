import Transaction from "@/views/transaction/Data";
import TransactionCreate from "@/views/transaction/Create";
import TransactionCreateCSV from "@/views/transaction/CreateFile";
import TransactionEdit from "@/views/transaction/Edit";

const transactionRouter = [
    { path: "/transaction", name: "transaction", component: Transaction },
    { path: "/transaction/create", name: "transaction.create", component: TransactionCreate },
    { path: "/transaction/create/csv", name: "transaction.create.csv", component: TransactionCreateCSV },
    { path: "/transaction/:no_transaction/edit", name: "transaction.edit", component: TransactionEdit },
]

export default transactionRouter;