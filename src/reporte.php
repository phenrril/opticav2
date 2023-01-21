<?php include_once "includes/header.php";
include "../conexion.php";
$id_user = $_SESSION['idUser'];
$permiso = "reporte";
$sql = mysqli_query($conexion, "SELECT p.*, d.* FROM permisos p INNER JOIN detalle_permisos d ON p.id = d.id_permiso WHERE d.id_usuario = $id_user AND p.nombre = '$permiso'");
$existe = mysqli_fetch_all($sql);
if (empty($existe) && $id_user != 1) {
    header("Location: permisos.php");
} ?>

<head>

    <link rel="stylesheet" href="https://use.fontawesome.com/releases/v5.8.1/css/all.css" integrity="sha384-50oBUHEmvpQ+1lW4y57PTFmhCaXp0ML5d60M1M7uH2+nqUivzIebhndOJK28anvf" crossorigin="anonymous"> 
    <script src='https://ajax.googleapis.com/ajax/libs/jquery/3.6.1/jquery.min.js'></script>
</head>
<br>

<div class="container is-fluid">
    <div class="col-xs-12">
        <h2 align="center">Reporte de Ventas</h2>
        <br><br>
        <div class="card">
        <div>
            <style>
                th {
                    font-weight: bold;
                    color: white;
                }
            </style>
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
                        <!-- <th>ID</th>
                        <th>Fecha</th>
                        <th>Descripcion</th>
                        <th>Nombre Cliente</th>
                        <th>metodo pago</th>
                        <th>Ingresos</th> -->
                        <th>ID USUARIO</th>
                        <th>ID VENTA</th>
                        <th>ID PRODUCTO</th>
                        <th>CANTIDAD</th>
                        <th>FECHA</th>
                        <th>PRECIO BRUTO</th>
                        <th>PRECIO NETO</th>
                        <th>TOTAL VENTA</th>
                        
                    </tr>
                </thead>
                <tbody>
                    <?php
                    $conexion = mysqli_connect("localhost", "root", "", "sis_venta");
                    if (isset($_GET['from_date']) && isset($_GET['to_date'])) {
                        $from_date = $_GET['from_date'];
                        $to_date = $_GET['to_date'];
                        //$query2 = mysqli_query("SELECT sum(ingresos) as 'subtotal' FROM  WHERE fecha BETWEEN '$from_date' AND '$to_date'");
                        //$query2 = mysqli_query($conexion, "SELECT sum(ingresos) as 'subtotal' FROM ingresos WHERE fecha BETWEEN '$from_date' AND '$to_date'"); 
                        
                        //$subtt = mysqli_fetch_assoc($query2);
                        // $query = "SELECT ingresos.*, cliente.nombre FROM ingresos
                        // JOIN cliente ON ingresos.id_cliente = cliente.idcliente
                        // WHERE ingresos.fecha BETWEEN '$from_date' AND '$to_date'";

                        // $query = "SELECT ingresos.*, cliente.nombre as 'nombre_cliente', metodos.descripcion as 'descripcion' FROM ingresos
                        // JOIN metodos ON ingresos.id_metodo = metodos.id
                        // JOIN cliente ON ingresos.id_cliente = cliente.idcliente 
                        // JOIN egresos on ingresos.id_cliente = egresos.id_cliente
                        // WHERE ingresos.fecha BETWEEN '$from_date' AND '$to_date'";


                        $query = "SELECT    detalle_venta.id_producto as 'idprod', 
                                            detalle_venta.cantidad as 'cantidad',
                                            detalle_venta.id_venta as 'idventa',
                                            ventas.id, ventas.total, ventas.id_usuario, ventas.fecha,
                                            producto.codproducto as 'id_prod', 
                                            producto.precio_bruto as 'preciobruto', 
                                            producto.precio as 'precioneto' from detalle_venta
                                            join ventas on detalle_venta.id_venta = ventas.id
                                            join producto on detalle_venta.id_producto = producto.codproducto
                                            WHERE ventas.fecha between '$from_date' AND '$to_date'";            
                
                        $totalventa = mysqli_query($conexion, "SELECT sum(total) as 'total' FROM ventas WHERE fecha BETWEEN '$from_date' AND '$to_date'");
                        $total = mysqli_fetch_assoc($totalventa);
                        $totalventab =  mysqli_query($conexion, "SELECT sum(producto.precio_bruto) as 'bruto', 
                                                                        detalle_venta.cantidad as 'cantidad',
                                                                        detalle_venta.id_venta as 'idventa',
                                                                        ventas.id, ventas.fecha,
                                                                        producto.codproducto as 'id_prod',
                                                                        producto.precio_bruto as 'preciobruto',
                                                                        sum(producto.precio) as 'precioneto' from detalle_venta
                                                                        join ventas on detalle_venta.id_venta = ventas.id
                                                                        join producto on detalle_venta.id_producto = producto.codproducto
                                                                        WHERE ventas.fecha between '$from_date' AND '$to_date'");
                        $totalb = mysqli_fetch_assoc($totalventab);
                        $totalventabruta = $totalb['bruto'];
                        $totalventaneta = $totalb['precioneto'];
                        $ganancia = $totalventaneta - $totalventabruta;


                        //$query="SELECT * FROM ingresos WHERE fecha BETWEEN '$from_date' AND '$to_date'";
                        
                        $query_run = mysqli_query($conexion, $query);
                        if (mysqli_num_rows($query_run) > 0) {
                            foreach ($query_run as $fila) {
                    ?>
                                <tr>
                                    <td><?php echo $fila['id_usuario']; ?></td>
                                    <td><?php echo $fila['idventa']; ?></td>
                                    <td><?php echo $fila['id_prod']; ?></td>
                                    <td><?php echo $fila['cantidad']; ?></td>
                                    <td><?php echo $fila['fecha']; ?></td>
                                    <td><?php echo $fila['preciobruto']; ?></td>
                                    <td><?php echo $fila['precioneto']; ?></td>
                                    <td><?php echo $fila['total']; ?></td>
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
                </tbody><td></td><td></td><td></td><td></td><td></td>
                <td><b>Total Venta Bruta: $<?php echo $totalventabruta;
                echo "<br>";
                ?>
                </b></td>
                <td><b>Total Venta Neta: $<?php echo $totalventaneta;
                echo "<br>";
                ?>
                </b></td>
                <td><b>Ganancia: $<?php echo $ganancia;
                echo "<br>";
                ?>
                </b></td>
            </table>
        </div>

    </div>
</div>
</div>
<script>      
$('#agregar_saldos').click(function () {
    var valor = document.getElementById('valor');
    if(valor.value == "" || valor.value == 0){   
    swal.fire
    ({
        position: 'top-end',
        showConfirmButton: false,
        title: 'Error',
        text: 'El valor no puede ser 0',
        icon: 'error'
    })
}
    else{
    if(confirm('¿Está seguro de agregar el valor? (no se puede cancelar)'))
    {
                {   
                $.ajax({
                        url: "saldos.php",
                        type: "POST",
                        data: $("#form_saldos").serialize(),
                        success: function (resultado){
                        $("#div_saldos").html(resultado);
                }
                });
        }
    }
}
})
            
</script>
<?php include_once "includes/footer.php"; ?>