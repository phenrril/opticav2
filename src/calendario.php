<?php include_once "includes/header.php";
include "../conexion.php";
$id_user = $_SESSION['idUser'];
$permiso = "calendario";
$sql = mysqli_query($conexion, "SELECT p.*, d.* FROM permisos p INNER JOIN detalle_permisos d ON p.id = d.id_permiso WHERE d.id_usuario = $id_user AND p.nombre = '$permiso'");
$existe = mysqli_fetch_all($sql);
if (empty($existe) && $id_user != 1) {
    header("Location: permisos.php");
} ?>

<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="stylesheet" href="https://use.fontawesome.com/releases/v5.8.1/css/all.css" integrity="sha384-50oBUHEmvpQ+1lW4y57PTFmhCaXp0ML5d60M1M7uH2+nqUivzIebhndOJK28anvf" crossorigin="anonymous">
    <link rel="stylesheet" href="../css/es.css">
    <title>Usuarios</title>
    <script src='https://ajax.googleapis.com/ajax/libs/jquery/3.6.1/jquery.min.js%27%3E'></script>
</head>
<br>

<div class="container is-fluid">
    <div class="col-xs-12">
        <h2>Calendario de Ventas, Ingresos y Egresos</h2>
        <br><br>
        <div class="card">
        <div>
            <style>
                th {
                    font-weight: bold;
                    color: white;
                }
            </style>
           
            <br>
            <br><br>


            <form action="" method="GET">
                <div class="row">
                    <div class="col-md-4">
                        <div class="form-group">
                            <label><b>Del Dia</b></label>
                            <input type="date" name="from_date" value="<?php if (isset($_GET['from_date'])) {
                                                                            echo $_GET['from_date'];
                                                                        } ?>" class="form-control">
                        </div>
                    </div>
                    <div class="col-md-4">
                        <div class="form-group">
                            <label><b> Hasta el Dia</b></label>
                            <input type="date" name="to_date" value="<?php if (isset($_GET['to_date'])) {
                                                                            echo $_GET['to_date'];
                                                                        } ?>" class="form-control">
                        </div>
                    </div>
                    <div class="col-md-4">
                        <div class="form-group">
                            <label><b></b></label> <br>
                            <button type="submit" class="btn btn-primary">Buscar</button>
                        </div>
                    </div>
                </div>
                <br>
            </form>
            <table class="table table-striped" id="table_id">
                <thead>
                    <tr class="bg-dark">
                        <th>ID Venta</th>
                        <th>ID Cliente</th>
                        <th>Total</th>
                        <th>Abonó</th>
                        <th>Restan</th>
                        <th>Fecha</th>
                    </tr>
                </thead>
                <tbody>
                    <?php
                    $conexion = mysqli_connect("localhost", "root", "", "sis_venta");
                    if (isset($_GET['from_date']) && isset($_GET['to_date'])) {
                        $from_date = $_GET['from_date'];
                        $to_date = $_GET['to_date'];
                        $query = "SELECT * FROM ventas WHERE fecha BETWEEN '$from_date' AND '$to_date'";
                        $query_run = mysqli_query($conexion, $query);
                        if (mysqli_num_rows($query_run) > 0) {
                            foreach ($query_run as $fila) {
                    ?>
                                <tr>
                                    <td><?php echo $fila['id']; ?></td>
                                    <td><?php echo $fila['id_cliente']; ?></td>
                                    <td><?php echo $fila['total']; ?></td>
                                    <td><?php echo $fila['abona']; ?></td>
                                    <td><?php echo $fila['resto']; ?></td>
                                    <td><?php echo $fila['fecha']; ?></td>
                                </tr>
                            <?php
                            }
                        } else {
                            ?>
                            <tr>
                                <td><?php echo "No se encontraron resultados"; ?></td>
                        <?php
                        
                        }
                    }
                        ?>
                        </tr>
                </tbody>
            </table>
        </div>
    </div>
</div>
</div>
  <script>      
//   document.querySelector("#agregar").addEventListener("click", function () {
//     {
//         $.ajax({
//             url: "saldos.php",
//             type: "POST",
//             data: $("#form_saldos").serialize(),
//             success: function (resultado) {
//                 $("#div_saldos").html(resultado);

//             }
//         });
//     }
// })
</script>
    